package p2p

import (
	"sync"

	"github.com/chronostech-git/fabrik/internal/log"
)

type PeerManager struct {
	Peers map[string]*Peer
	lock  sync.Mutex
}

// Create a new peer manager
func NewPeerManager() *PeerManager {
	return &PeerManager{
		Peers: make(map[string]*Peer),
	}
}

// Add peer to peer manager
func (m *PeerManager) AddPeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Peers[p.ID] = p
}

// Remove peer from peer manager
func (m *PeerManager) RemovePeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Peers, p.ID)
}

// Retrieve a peer from peer manager with an ID
func (m *PeerManager) GetPeer(ID string) *Peer {
	p, ok := m.Peers[ID]
	if !ok {
		log.Error("Peer with ID %s does not exist", ID)
		return nil
	}
	return p
}

// Broadcast will send a message to all peers (- sender) in the network.
// This is useful for blocks, and transactions, as all peers should be aware.
func (p *PeerManager) Broadcast(sender *Peer, msg *Message) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, peer := range p.Peers {
		if peer == sender {
			continue
		}
		peer.Send(msg)
	}
}
