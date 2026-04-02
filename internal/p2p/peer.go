package p2p

import (
	"bufio"
	"fmt"
	"net"
)

// A peer represents a single connection to another node
type Peer struct {
	ID     string
	Conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewInboundPeer(conn net.Conn) *Peer {
	return &Peer{
		ID:     conn.RemoteAddr().String(),
		Conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func NewOutboundPeer(network, host, port string) (*Peer, error) {
	address := net.JoinHostPort(host, port)
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &Peer{
		ID:     conn.RemoteAddr().String(), // remote node
		Conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (p *Peer) Send(msg *Message) error {
	_, err := fmt.Fprintf(p.writer, "%s %s", msg.Type, msg.Data)
	if err != nil {
		return err
	}
	return p.writer.Flush()
}
