package p2p

import (
	"sync"

	"github.com/chronostech-git/fabrik/internal/log"
)

type PeerManager struct {
	Peers map[string]*Peer
	lock  sync.Mutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		Peers: make(map[string]*Peer),
	}
}

func (m *PeerManager) AddPeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Peers[p.ID] = p
}

func (m *PeerManager) RemovePeer(p *Peer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Peers, p.ID)
}

func (m *PeerManager) GetPeer(ID string) *Peer {
	p, ok := m.Peers[ID]
	if !ok {
		log.Error("Peer with ID %s does not exist", ID)
		return nil
	}
	return p
}

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
