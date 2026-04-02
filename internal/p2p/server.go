package p2p

import (
	"log"
	"net"
)

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
