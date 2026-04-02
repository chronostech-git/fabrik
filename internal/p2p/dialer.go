package p2p

import (
	"strings"
)

// Dial peer connects to given address
// If connection is successful, the peer is added to the PeerManager
func DialPeer(address string, manager *PeerManager) error {
	host := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]

	peer, err := NewOutboundPeer("tcp", host, port)
	if err != nil {
		return err
	}

	manager.AddPeer(peer)

	go HandlePeer(peer, manager)

	return nil
}
