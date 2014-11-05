package main

import (
	"fmt"
)

type hub struct {
	// Inbound messages from the connections.
	broadcast chan *Message

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	rooms map[string]map[string]*connection
}

var h = hub{
	broadcast:  make(chan *Message),
	register:   make(chan *connection),
	unregister: make(chan *connection),
	//connections: make(map[*connection]bool),
	rooms: make(map[string]map[string]*connection),
}

// The hub registers connections by adding the connection pointer as a key
// in the connections map. The map value is always true. The hub unregisters
// connections by deleting the connection pointer from the connections map and
// closing the connection's send channel to signal the connection that no more
// messages will be sent to the connection.
//
// The hub handles messages by looping over the registered connections and sending
// the message to the connection's send channel. If the connection's send buffer
// is full, then the hub assumes that the client is dead or stuck. In this case,
// the hub unregisters the connection and closes the websocket.
func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			room_name := c.room.GetName()

			if _, ok := h.rooms[room_name]; ok {
			} else {
				h.rooms[room_name] = make(map[string]*connection)
			}

			h.rooms[room_name][c.user.Name] = c

			// @TODO change this to build response and send outside loop
			for _, conn := range h.rooms[room_name] {
				message := &Message{Action: "join", User: conn.user}
				msg := message.GetJSON()

				// sending "online users" to new connection
				c.send <- msg

				// sending new user to "online users"
				if conn.user.Name != c.user.Name {
					msg = fmt.Sprintf("{\"action\":\"join\",\"user\":{\"name\":\"%s\",\"headshot\":\"%d\"}}", c.user.Name, c.user.Headshot)
					conn.send <- msg
				}
			}

			messages := c.room.GetLastMessages(50)
			i := 49
			for _ = range messages {
				if nil != messages[i] {
					c.send <- messages[i].GetJSON()
				}

				i--
			}
		case c := <-h.unregister:
			room_name := c.room.GetName()
			msg := fmt.Sprintf("{\"action\":\"logout\",\"user\":{\"name\":\"%s\"}}", c.user.Name)
			delete(h.rooms[room_name], c.user.Name)
			close(c.send)

			// send user logout to connections
			for _, conn := range h.rooms[room_name] {
				conn.send <- msg
			}
		case message := <-h.broadcast:
			room_name := message.Room.GetName()

			for _, c := range h.rooms[room_name] {
				if room_name == c.room.GetName() {
					select {
					case c.send <- message.GetJSON():
					default:
						msg := fmt.Sprintf("{\"action\":\"logout\",\"user\":{\"name\":\"%s\"}}", c.user.Name)

						// send user logout to connections
						for _, conn := range h.rooms[room_name] {
							conn.send <- msg
						}

						delete(h.rooms[room_name], c.user.Name)
						close(c.send)
						go c.ws.Close()
					}
				}
			}
		}
	}
}
