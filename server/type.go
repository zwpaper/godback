package server

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
	ID     string `json:"id"`
	Number uint   `json:"number"`
	Err    string `json:"err"`
}

type roomEnterRequset struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	RoomID string `json:"room_id"`
}

type roomEnterResponse struct {
	ID     string `json:"id"`
	Number uint   `json:"number"`
	Err    string `json:"err"`
}
