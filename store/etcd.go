package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/client"
	"github.com/zwpaper/godback/utils"
	"golang.org/x/net/context"
)

// Client for etcd, Global
var (
	Client  client.Client
	prefix  string
	log     *logs.BeeLogger
	errInfo string
)

func init() {
	log = utils.Log
}

// Init use bindaddr in config to init etcd
// set gobal var EtcdClient
func Init(bindaddr []string, p string) (err error) {
	cfg := client.Config{
		Endpoints: bindaddr,
		Transport: client.DefaultTransport,
	}

	Client, err = client.New(cfg)
	if err != nil {
		return err
	}

	prefix = p
	op := &client.SetOptions{
		Dir: true}
	roomPath := path.Join(utils.PathRoom, utils.PathUsed)
	if _, err = get(roomPath); err != nil {
		err = set(roomPath, "", op)
		if err != nil {
			log.Emergency("Can not set room path in etcd")
			return err
		}
	}
	poolPath := path.Join(utils.PathRoom, utils.PathPool)
	pool := "0-999999"
	if _, err = get(poolPath); err != nil {
		err = set(poolPath, pool, nil)
		if err != nil {
			log.Emergency("Can not set room pool in etcd")
			return err
		}
	}

	return nil
}

// Get input etcd key, return client node
func get(key string) (*client.Node, error) {
	keyAPI := client.NewKeysAPI(Client)
	resp, err := keyAPI.Get(context.Background(), path.Join(prefix, key), nil)
	if err != nil {
		return nil, err
	}

	log.Debug("Got etcd nodes: %v", resp)
	return resp.Node, nil
}
func getFullPath(key string) (*client.Node, error) {
	keyAPI := client.NewKeysAPI(Client)
	resp, err := keyAPI.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}

	log.Debug("Got etcd nodes: %v", resp)
	return resp.Node, nil
}

func isKeyEmpty(key string) bool {
	eNode, err := get(key)
	if err != nil {
		errInfo := fmt.Sprintf("Key not found: %v\n%v", key, err)
		log.Error(errInfo)
		return true
	}
	if len(eNode.Nodes) > 1 {
		errInfo := fmt.Sprintf("Key %v is not empty", key)
		log.Error(errInfo)
		return false
	}
	return true
}

func set(key, val string, op *client.SetOptions) (err error) {
	keyAPI := client.NewKeysAPI(Client)
	_, err = keyAPI.Set(context.Background(), path.Join(prefix, key), val, op)
	if err != nil {
		log.Error("Set etcd error: %v", err.Error())
		return err
	}

	log.Debug("Set value %v to %v", val, key)
	return nil
}

func del(key string, op *client.DeleteOptions) (err error) {
	keyAPI := client.NewKeysAPI(Client)
	_, err = keyAPI.Delete(context.Background(), key, op)
	if err != nil {
		log.Error("Delete etcd error: %v", err.Error())
		return err
	}

	log.Info("Del from %v", key)
	return nil
}

// Room
func CreateRoom(room *Room) (string, error) {
	var errInfo string
	if room == nil {
		errInfo = fmt.Sprintf("No room sent in to create")
		log.Emergency(errInfo)
		return "", fmt.Errorf(errInfo)
	}

	ids, err := getAvailRoom()
	if err != nil {
		errInfo = fmt.Sprintf("Can not get available room in etcd: %v", err)
		log.Error(errInfo)
		return "", fmt.Errorf(errInfo)
	}

	op := &client.SetOptions{
		PrevExist: client.PrevNoExist}
	body, err := json.Marshal(room)
	if err != nil {
		errInfo := fmt.Sprintf("Can not marshal room request for etcd: %v", err)
		log.Error(errInfo)
		return "", fmt.Errorf(errInfo)
	}
	var id string
	for i := 0; i < utils.TimesRetry; i++ {
		id = strconv.Itoa(ids[rand.Intn(len(ids))])
		err = set(path.Join(utils.PathRoom, utils.PathUsed, id, utils.PathConfig),
			string(body), op)
		if err == nil {
			room.ID = id
			log.Info("Got room: %v", id)
			return id, nil
		}
	}

	errInfo = fmt.Sprintf("Can not found a room in etcd")
	log.Emergency(errInfo)
	return "", fmt.Errorf(errInfo)
}

func getAvailRoom() ([]int, error) {
	var errInfo string
	url := path.Join(utils.PathRoom, utils.PathUsed)
	usedRoomNode, err := get(url)
	if err != nil {
		errInfo = fmt.Sprintf("The used room dir no exist, creating...")
		log.Warn(errInfo)
		op := &client.SetOptions{
			Dir:       true,
			PrevExist: client.PrevNoExist}
		err = set(url, "", op)
		if err != nil {
			errInfo = fmt.Sprintf("Can not set the uesd room dir: %v",
				err)
			log.Emergency(errInfo)
			return nil, fmt.Errorf(errInfo)
		}

		usedRoomNode, err = get(url)
		if err != nil {
			errInfo = fmt.Sprintf(
				"Can not get used room after new create used room dir: %v", err)
			log.Emergency(errInfo)
			return nil, fmt.Errorf(errInfo)
		}
	}
	usedRoomMap := map[int]struct{}{}
	for _, n := range usedRoomNode.Nodes {
		id, err := strconv.Atoi(path.Base(n.Key))
		if err != nil {
			errInfo = fmt.Sprintf(
				"Room ID %v parsing error: %v", path.Base(n.Key), err)
			log.Emergency(errInfo)
			return nil, fmt.Errorf(errInfo)
		}
		usedRoomMap[id] = struct{}{}
	}

	poolNode, err := get(path.Join(utils.PathRoom, utils.PathPool))
	if err != nil {
		errInfo = fmt.Sprintf("Can not get Room pool: %v", err)
		log.Emergency(errInfo)
		return nil, fmt.Errorf(errInfo)
	}
	pool := strings.Split(poolNode.Value, "-")
	min, err := strconv.Atoi(pool[0])
	if err != nil {
		errInfo = fmt.Sprintf("Can not parse pool min: %v", pool)
		log.Emergency(errInfo)
		return nil, fmt.Errorf(errInfo)
	}
	max, err := strconv.Atoi(pool[1])
	if err != nil {
		errInfo = fmt.Sprintf("Can not parse pool max: %v", pool)
		log.Emergency(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	idInPool := make([]int, 0)
	for id := min; id <= max; id++ {
		if _, ok := usedRoomMap[id]; !ok {
			idInPool = append(idInPool, id)
			if len(idInPool) > utils.PoolSize {
				log.Info("Got id pool: %v", idInPool)
				return idInPool, nil
			}
		}
	}

	log.Info("Got id pool, not full", idInPool)
	return idInPool, nil
}

func GetRoom(id string) (*Room, error) {
	room := &Room{}
	url := path.Join(utils.PathRoom, utils.PathUsed, id, utils.PathConfig)
	roomNode, err := get(url)
	if err != nil {
		errInfo = fmt.Sprintf("Can not get room %v: %v", id, err)
		log.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	err = json.Unmarshal([]byte(roomNode.Value), room)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("Got room %v", room)
	return room, nil
}

func AddPlayerToRoom(roomID string, player *Player) (err error) {
	roomURL := path.Join(utils.PathRoom, utils.PathUsed, roomID)
	playerURL := path.Join(roomURL, utils.PathPlayer)
	_, err = get(roomURL)
	if err != nil {
		errInfo = fmt.Sprintf("Room %v not exist!", roomID)
		log.Error(errInfo)
		return fmt.Errorf(errInfo)
	}

	playersNode, err := get(playerURL)
	if err != nil {
		log.Notice("Player dir for room %v not exist, creating", roomID)
		op := &client.SetOptions{
			Dir: true}
		err = set(playerURL, "", op)
		if err != nil {
			errInfo = fmt.Sprintf("Can not create player dir for room %v", roomID)
			log.Emergency(errInfo)
			return fmt.Errorf(errInfo)
		}

		playersNode, err = get(playerURL)
		if err != nil {
			errInfo = fmt.Sprintf(
				"Can not get players after creating player dir for room %v", roomID)
			log.Emergency(errInfo)
			return fmt.Errorf(errInfo)
		}
	}

	for i := len(playersNode.Nodes) + 1; i < 16; i++ {
		player.Order = i
		playerJSON, err := json.Marshal(*player)
		if err != nil {
			log.Emergency(err.Error())
			return err
		}

		err = set(path.Join(playerURL, strconv.Itoa(i)), string(playerJSON), nil)
		if err != nil {
			errInfo = fmt.Sprintf("Can not add player %v to room %v: %v",
				player.Name, roomID, err)
			log.Notice(errInfo)
			continue
		}
		return nil
	}
	return fmt.Errorf("Can not add player %v to room %v", player, roomID)
}

func GetAllPlayersInRoom(roomID string) (*[]Player, error) {
	playerURL := path.Join(utils.PathRoom, utils.PathUsed, roomID, utils.PathPlayer)
	playersNode, err := get(playerURL)
	if err != nil {
		errInfo = fmt.Sprintf("Can not get player in %v!", roomID)
		log.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	players := &[]Player{}
	player := &Player{}
	for _, node := range playersNode.Nodes {
		playerKey := node.Key
		playerNode, err := getFullPath(playerKey)
		if err != nil {
			errInfo = fmt.Sprintf("Can not get player %v config: %v",
				playerKey, err)
			log.Error(errInfo)
			return nil, fmt.Errorf(errInfo)
		}
		err = json.Unmarshal([]byte(playerNode.Value), player)
		if err != nil {
			errInfo = fmt.Sprintf("Can not get player %v config: %v",
				playerKey, err)
			log.Error(errInfo)
			return nil, fmt.Errorf(errInfo)
		}

		*players = append(*players, *player)
	}
	log.Info("Get players: %v", players)
	return players, nil
}

func GetPlayerInRoom(playerID, roomID string) (*Player, error) {
	players, err := GetAllPlayersInRoom(roomID)
	if err != nil {
		return nil, err
	}

	for _, p := range *players {
		if p.ID == playerID {
			log.Info("Got player %v", p)
			return &p, nil
		}
	}

	errInfo = fmt.Sprintf("No player %v found", playerID)
	log.Error(errInfo)
	return nil, fmt.Errorf(errInfo)
}

func UpdatePlayerInRoom(roomID string, player *Player) (err error) {
	roomURL := path.Join(utils.PathRoom, utils.PathUsed, roomID)
	playerURL := path.Join(roomURL, utils.PathPlayer, strconv.Itoa(player.Order))

	_, err = get(playerURL)
	if err != nil {
		errInfo = fmt.Sprintf(
			"Can not get player %v in room %v", player, roomID)
		log.Emergency(errInfo)
		return fmt.Errorf(errInfo)
	}

	playerJSON, err := json.Marshal(*player)
	if err != nil {
		log.Emergency(err.Error())
		return err
	}

	err = set(playerURL, string(playerJSON), nil)
	if err != nil {
		errInfo = fmt.Sprintf("Can not update player %v to room %v: %v",
			player.Name, roomID, err)
		log.Notice(errInfo)
		return fmt.Errorf(errInfo)
	}
	return nil
}
