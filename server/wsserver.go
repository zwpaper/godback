package server

type Handler func() string

type Game struct {
	// device id to client
	clients map[string]*Client

	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	size uint
	// State Machine
	Handlers   map[string]Handler
	StartState string
	EndStates  map[string]bool
	Pipe       chan *gameRequest
	End        chan struct{}
}

func (g *Game) run() {
	for {
		select {
		case client := <-g.register:
			logger.Info("%v resigter to game", client.ID)
			g.clients[client.ID] = client
		case client := <-g.unregister:
			if _, ok := g.clients[client.ID]; ok {
				delete(g.clients, client.ID)
				close(client.send)
			}
		case message := <-g.broadcast:
			for _, client := range g.clients {
				select {
				case client.send <- message:
					//				default:
					//					close(client.send)
					//					delete(g.clients, client.ID)
				}
			}
		}
	}
}
