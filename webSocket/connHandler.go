package ws

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"Practice-/tracker"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func generatePeerID() string {
	return fmt.Sprintf("peer-%d", rand.Intn(1000000))
}

func HandleWebSocket(t *tracker.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Connection error:", err)
			return
		}
		defer conn.Close()

		peerID := generatePeerID()
		t.AddPeer(peerID)
		log.Println("New peer connected:", peerID)
		conn.WriteMessage(websocket.TextMessage, []byte("Your peer ID: "+peerID))

		// Send current peer list to the new peer
		peerList := t.ListPeers()
		conn.WriteMessage(websocket.TextMessage, []byte("Connected peers: "+fmt.Sprint(peerList)))

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Peer disconnected:", peerID)
				t.RemovePeer(peerID)
				break
			}
		}
	}
}
