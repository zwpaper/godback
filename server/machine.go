package server

import (
	"encoding/json"
	"fmt"

	"github.com/zwpaper/godback/store"
)

func (g *Game) AddState(handlerName string, handlerFn Handler) {
	g.Handlers[handlerName] = handlerFn
}

func (g *Game) AddEndState(endState string) {
	g.EndStates[endState] = true
}

func (g *Game) Execute() {
	logger.Info("State machine start!")
	if handler, present := g.Handlers[g.StartState]; present {
		for {
			nextState := handler()
			_, finished := g.EndStates[nextState]
			if finished {
				break
			} else {
				handler, present = g.Handlers[nextState]
			}
		}
	}
}

const (
	stateEnter = "enter"
	stateReady = "ready"
	stateEnd   = "end"
)

func newGame(n uint) (g *Game) {
	g = &Game{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		size:       n,
		Handlers:   make(map[string]Handler),
		StartState: stateEnter,
		EndStates:  make(map[string]bool),
		Pipe:       make(chan *gameRequest),
	}
	g.AddState(stateEnter, g.enterGame)

	return g
}

func (g *Game) enterGame() string {
	logger.Info("Enter state: %v", stateEnter)
	for {
		select {
		case request := <-g.Pipe:
			if request.OP != stateEnter {
				errInfo = fmt.Sprintf("Request not match state %v", stateEnter)
				logger.Emergency(errInfo)
				continue
			}

			response := &gameResponse{
				OP:      stateEnter,
				Success: false}
			players := &[]store.Player{}
			room, err := store.GetRoom(request.RoomID)
			if err != nil {
				errInfo = fmt.Sprintf("Can not get room %v info!\n%v",
					request.RoomID, err)
				logger.Error(errInfo)
				continue
			}
			err = store.AddPlayerToRoom(request.RoomID, &store.Player{
				ID:   request.UID,
				Name: request.Name})
			if err != nil {
				errInfo = fmt.Sprintf("Can not add room creater to room!\n%v", err)
				logger.Error(errInfo)
				continue
			}

			players, err = store.GetAllPlayersInRoom(request.RoomID)
			if err != nil {
				errInfo = fmt.Sprintf(
					"Can not get players in room %v \n%v", request.RoomID, err)
				logger.Error(errInfo)
				continue
			}
			logger.Info("Added %v to room %v", request.Name, request.RoomID)
			response.Players = *players
			response.Success = true
			response.Err = ""
			response.Number = countPlayers(room)
			logger.Info("Response: %v", response)

			msg, err := json.Marshal(response)
			if err != nil {
				logger.Emergency(err.Error())
				continue
			}
			g.broadcast <- msg
			if len(g.clients) == int(g.size) {
				logger.Info("All player in, goto", stateReady)
				return stateReady
			}
		case <-g.End:
			return stateEnd
		}
	}
}
