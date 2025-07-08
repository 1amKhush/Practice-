package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	ma "github.com/multiformats/go-multiaddr"
)

const TrackerProtocol = "/tracker/1.0.0"

func main() {
	ctx := context.Background()
	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[CLIENT] My Peer ID:", h.ID())

	// -- STEP 1: Ask for tracker address
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter tracker multiaddr: ")
	addrStr, _ := reader.ReadString('\n')
	addrStr = strings.TrimSpace(addrStr)

	// Example: /ip4/127.0.0.1/tcp/9000/p2p/QmTrackerID
	maddr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		log.Fatal("Invalid multiaddr:", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal("Failed to get peer info:", err)
	}

	// -- STEP 2: Connect to tracker
	if err := h.Connect(ctx, *peerInfo); err != nil {
		log.Fatal("Connection failed:", err)
	}
	fmt.Println("[CLIENT] Connected to tracker.")

	// -- STEP 3: Open a stream
	s, err := h.NewStream(ctx, peerInfo.ID, protocol.ID(TrackerProtocol))
	if err != nil {
		log.Fatal("Stream error:", err)
	}

	// -- STEP 4: Send peer name
	peerName := fmt.Sprintln("Khushvendra")

	msg := fmt.Sprintf("%s\n", peerName)
	_, _ = s.Write([]byte(msg))

	// -- STEP 5: Read response
	r := bufio.NewReader(s)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("[CLIENT] Disconnected from tracker.")
			break
		}
		fmt.Print("[TRACKER] ", line)
	}

	// STEP 6: Send periodic heartbeats
	// go func() {
	// 	for {
	// 		_, err := s.Write([]byte("ping\n"))
	// 		if err != nil {
	// 			fmt.Println("[CLIENT] Failed to send ping:", err)
	// 			return
	// 		}
	// 		fmt.Println("[CLIENT] Ping sent.")
	// 		time.Sleep(10 * time.Second) // send ping every 10 seconds
	// 	}
	// }()

	// STEP 7: Keep alive
	select {}
}
