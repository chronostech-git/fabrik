package p2p

import (
	"log"
	"sync"
)

type PeerManager struct {
	Peers map[string]*Peer
	lock  sync.Mutex
}

// NewPeerManager create a new peer manager
func NewPeerManager() *PeerManager {
	return &PeerManager{
		Peers: make(map[string]*Peer),
	}
}

// AddPeer add peer to peer manager
func (m *PeerManager) AddPeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Peers[p.ID] = p
}

// RemovePeer remove peer from peer manager
func (m *PeerManager) RemovePeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Peers, p.ID)
}

// GetPeer retrieve a peer from peer manager with an ID
func (m *PeerManager) GetPeer(ID string) *Peer {
	p, ok := m.Peers[ID]
	if !ok {
		log.Panicf("Peer with ID %s does not exist", ID)
		return nil
	}
	return p
}

// Broadcast sends a message to all peers (- sender) in the network.
// This is useful for blocks, and transactions, as all peers should be aware.
func (p *PeerManager) Broadcast(sender *Peer, msg *Message) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, peer := range p.Peers {
		if peer == sender {
			continue
		}
		err := peer.Send(msg)
		if err != nil {
			log.Println("Failed to send message to peer")
			break
		}
	}
}
