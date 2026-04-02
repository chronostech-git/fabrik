package p2p

import (
	"log"
	"net"
)

// StartServer is used in the cmd/fabnet folder.
// It starts the server that peers connect to to find/see one another.
// Later, this will be modified to somehow run automatically when using
// cli/node.
func StartServer(address string, manager *PeerManager) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		peer := NewInboundPeer(conn)
		manager.AddPeer(peer)

		go HandlePeer(peer, manager)
	}
}
