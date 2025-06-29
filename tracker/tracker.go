package tracker

import (
	"sync"
)

type Tracker struct {
	mu    sync.Mutex
	peers map[string]bool
}

func NewTracker() *Tracker {
	return &Tracker{
		peers: make(map[string]bool),
	}
}

func (t *Tracker) AddPeer(peerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peers[peerID] = true
}

func (t *Tracker) RemovePeer(peerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.peers, peerID)
}

func (t *Tracker) ListPeers() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	var list []string
	for peer := range t.peers {
		list = append(list, peer)
	}
	return list
}
