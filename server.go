package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

func main() {
	go h.run()

	http.HandleFunc("/", findHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))

	fmt.Println("Serving on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
