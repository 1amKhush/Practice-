// cmd/peer/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "github.com/1amKhush/Practice-/p2p"
	// "github.com/1amKhush/Practice-/tracker"
	"github.com/1amKhush/Practice-/webRTC"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

const TrackerProtocol = "/tracker/1.0.0"

func main() {
	// 1. Start libp2p host
	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[PEER] Host ID:", h.ID())

	// 2. Ask for tracker address (file or manual)
	addrStr := readTrackerAddr()

	// 3. Connect to tracker
	maddr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		log.Fatal("Invalid tracker address:", err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := h.Connect(context.Background(), *info); err != nil {
		log.Fatal("Failed to connect to tracker:", err)
	}
	fmt.Println("[PEER] Connected to tracker.")

	// 4. Open tracker stream to register and get peer list
	s, err := h.NewStream(context.Background(), info.ID, protocol.ID(TrackerProtocol))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// send our name and listen for list
	fmt.Print("Enter your name: ")
	name := bufio.NewReader(os.Stdin)
	peerName, _ := name.ReadString('\n')
	s.Write([]byte(strings.TrimSpace(peerName) + "\n"))

	reader := bufio.NewReader(s)
	// read welcome
	welcome, _ := reader.ReadString('\n')
	fmt.Print("[TRACKER] ", welcome)
	// read peers list
	peersLine, _ := reader.ReadString('\n')
	fmt.Print("[TRACKER] ", peersLine)

	// 5. Initialize WebRTC peer
	webrtcPeer, err := webRTC.NewWebRTCPeer()
	if err != nil {
		log.Fatal("WebRTC init error:", err)
	}
	defer webrtcPeer.Close()

	// CLI loop for WebRTC offers/answers/files
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)
		switch parts[0] {
		case "offer":
			handleOffer(webrtcPeer)
		case "answer":
			if len(parts) < 2 {
				fmt.Println("Usage: answer <offer_json>")
				continue
			}
			handleAnswer(webrtcPeer, strings.Join(parts[1:], " "))
		case "complete":
			if len(parts) < 2 {
				fmt.Println("Usage: complete <answer_json>")
				continue
			}
			handleComplete(webrtcPeer, strings.Join(parts[1:], " "))
		case "download":
			if len(parts) != 2 {
				fmt.Println("Usage: download <filename>")
				continue
			}
			err := webrtcPeer.RequestFile(parts[1])
			if err != nil {
				fmt.Println("Error requesting file:", err)
			}
		case "exit":
			return
		default:
			fmt.Println("Unknown command. Type help for instructions.")
		}
	}
}

func readTrackerAddr() string {
	// attempt to read from file
	if data, err := os.ReadFile("tracker.addr"); err == nil {
		return strings.TrimSpace(string(data))
	}
	// fallback to manual
	fmt.Print("Tracker address: ")
	addr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(addr)
}

func handleOffer(p *webRTC.WebRTCPeer) {
	offer, err := p.CreateOffer()
	if err != nil {
		fmt.Println("Offer error:", err)
		return
	}
	fmt.Println("Offer JSON:", offer)
}

func handleAnswer(p *webRTC.WebRTCPeer, offer string) {
	answer, err := p.CreateAnswer(offer)
	if err != nil {
		fmt.Println("Answer error:", err)
		return
	}
	fmt.Println("Answer JSON:", answer)
}

func handleComplete(p *webRTC.WebRTCPeer, ans string) {
	if err := p.SetAnswer(ans); err != nil {
		fmt.Println("Complete error:", err)
		return
	}
	fmt.Println("Waiting for connection...")
	if err := p.WaitForConnection(30 * time.Second); err != nil {
		fmt.Println("Connection timeout", err)
	}
}
