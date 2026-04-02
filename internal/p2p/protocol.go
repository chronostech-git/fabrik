package p2p

import (
	"bufio"

	"github.com/chronostech-git/fabrik/internal/log"
)

// HandlePeer handles messages, and filters them based on msg.Type.
// If msg.Type is a block or transaction, it will broadcast it to all known peers.
// If msg.Type is something like a handshake, it will be between only 2 peers.
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
