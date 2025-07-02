package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1amKhush/Practice-/db"
	"github.com/1amKhush/Practice-/tracker"
	"github.com/1amKhush/Practice-/webSocket"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to access .env file")
	}
	db.InitDB()

	t := tracker.NewTracker()

	http.HandleFunc("/ws", ws.HandleWebSocket(t))
	
	fmt.Println("-> WebSocket server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
