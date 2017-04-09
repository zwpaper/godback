package store

type Charactor struct {
	Wolves    uint
	Villagers uint
	Prophet   bool
	Witch     bool
	Hunter    bool
	KingWolf  bool
	Guard     bool
}

type Room struct {
	ID   string
	Char Charactor
}

type Player struct {
	ID     string `json:"id"`
	Order  int    `json:"order"`
	Name   string `json:"name"`
	Char   string `json:"char"`
	Status string `json:"status"`
}
