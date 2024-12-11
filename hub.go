package main

// Hub maintains the set of active clients and broadcasts messages to the
type Hub struct {
	// put registered clients into the room.
	rooms map[string]map[*connection]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan subscription

	// Unregister requests from clients.
	unregister chan subscription
}

type message struct {
	Room string `json:"room"`
	Data []byte `json:"data"`
}

var H = &Hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
	// string is like a room id, the value is who connects to it
	// map[string]map[*connection]bool and not map[string]*connection to handle multiple users per room
}

func (h *Hub) Run() {
	// wait for websocket event forever (like a while true)
	for {
		select {
		// if there's a register message from the channel
		case s := <-h.register:
			connections := h.rooms[s.room] // return map[*connection]bool
			// if there's no connection...
			if connections == nil {
				// ...make one
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			// set it to true
			h.rooms[s.room][s.conn] = true
			// if there's an unregister message from the channel
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		// if there's a broadcast message from the channel
		case m := <-h.broadcast:
			connections := h.rooms[m.Room]
			for c := range connections {
				select {
				case c.send <- m.Data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.Room)
					}
				}
			}
		}
	}
}
