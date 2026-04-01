package p2p

import (
	"fmt"

	"github.com/google/uuid"
)

func CreatePeerLink(address string, peerID uuid.UUID) string {
	return fmt.Sprintf("tcp:///%s/%s", address, peerID)
}
