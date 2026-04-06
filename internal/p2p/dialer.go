package p2p

import (
	"strings"
)

// DialPeer dials a peer given a peer address, and adds the peer
// to the PeerManager
func DialPeer(address string, manager *PeerManager) (*Peer, error) {
	host := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]

	peer, err := NewOutboundPeer("tcp", host, port)
	if err != nil {
		return nil, err
	}

	manager.AddPeer(peer)

	go HandlePeer(peer, manager)

	return peer, err
}
