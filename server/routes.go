package server

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/zwpaper/godinwerewolves/store"
)

var (
	log     *logs.BeeLogger
	errInfo string
)

func setRoute(r *gin.Engine) {
	r.POST("/room", createRoom)
	r.POST("/room/:room", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "room",
			"id":      c.Param("room")})
	})
	r.POST("/room/:room/player", enterRoom)
}

// Room
func createRoom(c *gin.Context) {
	log.Debug("Received createRoom request")
	var err error
	request := &roomCreationRequset{}
	response := &roomCreationResponse{}
	if err = c.BindJSON(request); err != nil {
		errInfo = fmt.Sprintf("Can not parse the request: %v", err)
		log.Error(errInfo)
		response.Err = errInfo
		c.JSON(http.StatusBadRequest, response)
		return
	}
	log.Debug("%v", request)

	room := &store.Room{
		Char: store.Charactor{
			Wolves:    request.Wolves,
			Villagers: request.Villagers,
			Prophet:   request.Prophet,
			Witch:     request.Witch,
			Hunter:    request.Hunter,
			KingWolf:  request.KingWolf,
			Guard:     request.Guard}}
	id, err := store.CreateRoom(room)
	if err != nil {
		errInfo := fmt.Sprintf("Can not create room: %v", err)
		response.Err = errInfo
		log.Error(errInfo)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	err = store.AddPlayerToRoom(id, &store.Player{
		ID:   request.UID,
		Name: request.Name})
	if err != nil {
		errInfo = fmt.Sprintf("Can not add room creater to room!\n%v", err)
		log.Emergency(errInfo)
		response.Err = errInfo
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.ID = id
	response.Err = ""
	response.Number = countPlayers(room)
	c.JSON(http.StatusCreated, response)
	return
}

func countPlayers(room *store.Room) uint {
	count := room.Char.Wolves + room.Char.Villagers
	if room.Char.Prophet {
		count++
	}
	if room.Char.Witch {
		count++
	}
	if room.Char.Hunter {
		count++
	}
	if room.Char.KingWolf {
		count++
	}
	if room.Char.Guard {
		count++
	}
	return count
}

func enterRoom(c *gin.Context) {
	log.Debug("Received enter room request")
	var err error
	request := &roomEnterRequset{}
	response := &roomCreationResponse{}
	if err = c.BindJSON(request); err != nil {
		errInfo = fmt.Sprintf("Can not parse the request: %v", err)
		log.Error(errInfo)
		response.Err = errInfo
		c.JSON(http.StatusBadRequest, response)
		return
	}
	log.Debug("%v", request)

	room, err := store.GetRoom(request.RoomID)
	if err != nil {
		errInfo = fmt.Sprintf("Can not get room %v info!\n%v", request.RoomID, err)
		log.Error(errInfo)
		response.Err = errInfo
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	err = store.AddPlayerToRoom(request.RoomID, &store.Player{
		ID:   request.UID,
		Name: request.Name})
	if err != nil {
		errInfo = fmt.Sprintf("Can not add room creater to room!\n%v", err)
		log.Emergency(errInfo)
		response.Err = errInfo
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.ID = request.RoomID
	response.Err = ""
	response.Number = countPlayers(room)
	c.JSON(http.StatusCreated, response)
	return
}
