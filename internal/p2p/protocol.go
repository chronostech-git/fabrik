package p2p

import (
	"bufio"

	"github.com/chronostech-git/fabrik/internal/log"
)

func HandlePeer(peer *Peer, manager *PeerManager) {
	scanner := bufio.NewScanner(peer.Conn)

	for scanner.Scan() {

		line := scanner.Text()

		msg, err := ParseMessage(line)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		switch msg.Type {

		case "PING":
			peer.Send(&Message{"PONG", ""})

		case "TX":
			manager.Broadcast(peer, msg)

		case "BLOCK":
			manager.Broadcast(peer, msg)

		}
	}

	manager.RemovePeer(peer)
	peer.Conn.Close()
}
