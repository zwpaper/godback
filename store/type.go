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
	ID     string
	Order  int
	Name   string
	Char   string
	Status string
}
