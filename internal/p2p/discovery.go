package p2p

import (
	"bufio"
	"fmt"
	"net"
)

type Handshake struct {
	instigator *Peer
	receiver   string
	disk       *DiskStorage
	msg        *Message
	manager    *PeerManager
	done       bool
}

// NewHandshake creates an instance of Handshake
func NewHandshake(
	instigator *Peer,
	receiver string,
	disk *DiskStorage,
	manager *PeerManager,
) *Handshake {
	return &Handshake{
		instigator: instigator,
		receiver:   receiver,
		disk:       disk,
		manager:    manager,
		done:       true,
	}
}

// Send executes the handshake by sending the message to the receiver.
func (h *Handshake) Send() error {
	// Dial the receiver directly
	conn, err := net.Dial("tcp", h.receiver)
	if err != nil {
		return fmt.Errorf("failed to connect to receiver %s: %w", h.receiver, err)
	}
	defer conn.Close() // close after sending

	// Buffered writer for the connection
	writer := bufio.NewWriter(conn)

	h.msg = &Message{
		Type: "HANDSHAKE",
		Data: fmt.Sprintf("Handshake between %s and %s successful", h.instigator.ID, h.receiver),
	}

	// Send the handshake message
	_, err = fmt.Fprintf(writer, "%s %s\n", h.msg.Type, h.msg.Data)
	if err != nil {
		return fmt.Errorf("failed to send handshake to %s: %w", h.receiver, err)
	}

	// Ensure the message is flushed
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush handshake to %s: %w", h.receiver, err)
	}

	return nil
}
