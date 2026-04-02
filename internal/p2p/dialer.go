package p2p

import (
	"strings"
)

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
