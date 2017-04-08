package server

import "github.com/zwpaper/godback/store"

type roomCreationRequset struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Wolves    uint   `json:"wolves"`
	Villagers uint   `json:"villagers"`
	Prophet   bool   `json:"prophet"`
	Witch     bool   `json:"witch"`
	Hunter    bool   `json:"hunter"`
	KingWolf  bool   `json:"kingwolf"`
	Guard     bool   `json:"guard"`
}

type roomCreationResponse struct {
	ID     string `json:"room_id"`
	Number uint   `json:"number"`
	Err    string `json:"err"`
}

type gameRequest struct {
	UID    string `json:"uid"`
	OP     string `json:"op"`
	Name   string `json:"name"`
	RoomID string `json:"room_id"`
}

type gameResponse struct {
	ID      string         `json:"id"`
	OP      string         `json:"op"`
	Number  uint           `json:"number"`
	Players []store.Player `json:"players"`
	Success bool           `json:"success"`
	Err     string         `json:"err"`
}
