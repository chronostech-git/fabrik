package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

// A peer represents a single connection to another node
type Peer struct {
	ID     string `json:"id"`
	Conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// Create a new "inbound" peer.
// An inbound peer uses a RemoteAddr (remote network address) if known.
func NewInboundPeer(conn net.Conn) *Peer {
	return &Peer{
		ID:     conn.RemoteAddr().String(),
		Conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

// Create an "outbound" peer.
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

// Send a message from specific peer.
// See message.go for msg implementation.
func (p *Peer) Send(msg *Message) error {
	_, err := fmt.Fprintf(p.writer, "%s\n", msg)
	if err != nil {
		return err
	}
	return p.writer.Flush()
}

// For disk file.
// This allows for json marshaling compatibility.
type peerJson struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// PeerToJson converts a given peer into a json string.
func PeerToJson(p *Peer) string {
	host := strings.Split(p.Conn.RemoteAddr().String(), ":")[0]
	port := strings.Split(p.Conn.RemoteAddr().String(), ":")[1]

	data := &peerJson{
		ID:   p.ID,
		Host: host,
		Port: port,
	}

	j, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	return string(j)
}
