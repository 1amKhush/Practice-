package ws

import (
	
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"

	"strings"

	
	"github.com/1amKhush/Practice-/tracker"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func generatePeerID() string {
	id := rand.Intn(1000000)
	return fmt.Sprint(id)
}

func peerIPAddr(r *http.Request) (string) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("Could not fetch IP of connected PEER")
	}

	return fmt.Sprint(host)
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
		peerIP := peerIPAddr(r)

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading peer name:", err)
			return
		}
		name := strings.TrimSpace(string(msg))

		if err := t.AddPeer(peerID, name, peerIP); err != nil {
    		log.Println("DB error saving peer:", err)
		}

		log.Printf("-> '%s' connected with ID: %s", name, peerID)
		conn.WriteMessage(websocket.TextMessage, []byte("Your peer ID: "+peerID))
		conn.WriteMessage(websocket.TextMessage, []byte(peerIP))

		peerList := t.ListPeers()
		conn.WriteMessage(websocket.TextMessage, []byte("Connected peers: "+fmt.Sprint(peerList)))

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("-> '%s' disconnected: %s", name, peerID)
				t.RemovePeer(peerID)
				break
			}
		}
	}
}
