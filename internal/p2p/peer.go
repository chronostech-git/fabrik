package p2p

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

type Peer struct {
	Link      string     // Example: tcp:///127.0.0.1:6000/<peer_id>
	EventFeed *EventFeed // EventFeed will handle communication between peers, for example broadcasting a block.
	Conn      net.Conn
}

func NewPeer(network string, address string) *Peer {
	peerAddr := fmt.Sprintf("%s", address)

	conn, err := net.Dial(network, peerAddr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	peerID, err := uuid.NewUUID()
	if err != nil {
		log.Panic(err)
	}

	peerLink := CreatePeerLink(address, peerID)

	return &Peer{
		Link:      peerLink,
		EventFeed: &EventFeed{},
		Conn:      conn,
	}
}
