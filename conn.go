package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string

	// User associated with connection.
	user *User

	// Name of room associated with connection.
	// rooms map[string]bool
	room *Room
}

// The reader method transfers inbound messages
// from the websocket to the hub.
func (c *connection) reader() {
	for {
		var msg string
		if err := websocket.Message.Receive(c.ws, &msg); err != nil {
			if err.Error() != "EOF" {
				fmt.Println("connection.reader(): ", err)
			}
			break
		}

		message, err := parse_message(msg)
		if err != nil {
			fmt.Println(err)
		}

		message.User = c.user

		if err == nil && message.Save() {
			h.broadcast <- message
		}
	}
	c.ws.Close()
}

// The writer method transfers messages from
// the connection's send channel to the websocket.
// The writer method exits when the channel is closed
// by the hub, or there's an error writing to the
// websocket.
func (c *connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
