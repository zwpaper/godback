package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/zwpaper/godback/store"
	"github.com/zwpaper/godback/utils"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// Device id
	ID string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	request := &gameRequset{}
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logs.Error("error: %v", err)
			}
			break
		}

		err = json.Unmarshal(message, request)
		if err != nil {
			errInfo = fmt.Sprintf("When receive game request: %v", err)
			logs.Emergency(errInfo)
			continue
		}
		switch {
		case request.OP == utils.OPEnter:
			c.handleEnterRoom(request)
		}
		//c.hub.broadcast <- message
	}
}

func (c *Client) handleEnterRoom(r *gameRequset) {
	response := &gameResponse{
		OP:      utils.OPEnterSucc,
		Success: false}
	var (
		players *[]store.Player
	)
	room, err := store.GetRoom(r.RoomID)
	if err != nil {
		errInfo = fmt.Sprintf("Can not get room %v info!\n%v", r.RoomID, err)
		goto ErrorReturn
	}
	err = store.AddPlayerToRoom(r.RoomID, &store.Player{
		ID:   r.UID,
		Name: r.Name})
	if err != nil {
		errInfo = fmt.Sprintf("Can not add room creater to room!\n%v", err)
		goto ErrorReturn
	}

	players, err = store.GetAllPlayersInRoom(r.RoomID)
	if err != nil {
		errInfo = fmt.Sprintf(
			"Can not get players in room %v \n%v", r.RoomID, err)
		goto ErrorReturn
	}
	logs.Info("Added %v to room %v", r.Name, r.RoomID)
	response.Players = *players
	response.Success = true
	response.Err = ""
	response.Number = countPlayers(room)
	logs.Info("Response: %v", response)
	c.conn.WriteJSON(response)
	return

ErrorReturn:
	logs.Error(errInfo)
	response.Err = errInfo
	c.conn.WriteJSON(response)
	return

}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("%v", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	logs.Info("added a client to hub")
	go client.writePump()
	client.readPump()
}
