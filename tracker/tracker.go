package tracker


type Tracker struct {
	peers map[string]bool
}

func NewTracker() *Tracker {
	return &Tracker{
		peers: make(map[string]bool),
	}
}

func (t *Tracker) AddPeer(peerID string) {
	t.peers[peerID] = true
}

func (t *Tracker) RemovePeer(peerID string) {
	delete(t.peers, peerID)
}

func (t *Tracker) ListPeers() []string {
	var list []string
	for peer := range t.peers {
		list = append(list, peer)
	}
	return list
}

func logNewPeer()