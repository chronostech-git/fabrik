package p2p

import (
	"log"
	"net"
)

// StartServer starts the networking server and listens for incoming
// peer connections.
func StartServer(address string, manager *PeerManager, disk *DiskStorage) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		log.Println("[FABNET] Peer connected with address", conn.RemoteAddr())

		peer := NewInboundPeer(conn)
		manager.AddPeer(peer)

		disk.WritePeer(peer)

		go HandlePeer(peer, manager)
	}
}
