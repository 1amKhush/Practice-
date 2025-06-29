package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1amKhush/Practice-/tracker"
	"github.com/1amKhush/Practice-/webSocket"
)

func main() {
	// Initialize peer tracker
	t := tracker.NewTracker()

	// WebSocket handler with tracker reference
	http.HandleFunc("/ws", ws.HandleWebSocket(t))

	fmt.Println("WebSocket server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
