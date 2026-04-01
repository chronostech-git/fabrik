package p2p

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var ErrUnknownPeer = errors.New("unknown peer")

type Server struct {
	Peers map[string]*Peer
	Lock  sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Peers: make(map[string]*Peer),
	}
}

func (s *Server) AddPeer(p *Peer) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Peers[p.Link] = p
	log.Println("Peer added", p.Link)
}

func (s *Server) RemovePeer(p *Peer) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Peers, p.Link)
	log.Println("Peer removed", p.Link)
}

func (s *Server) Broadcast(message *Message) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	for _, peer := range s.Peers {
		fmt.Fprintln(peer.Conn, message.PeerLink, message.Data)
	}
}

func (s *Server) handleIncomingConnection(conn net.Conn) {
	defer conn.Close()

	peer := NewPeer(conn.LocalAddr().Network(), conn.LocalAddr().String())
	s.AddPeer(peer)
	defer s.RemovePeer(peer)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("%s: %s", peer.Link, msg)
		payload := NewMessage(peer.Link, msg)
		s.Broadcast(payload)
	}
	log.Println("Connection closed")
}

func (s *Server) Start(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	log.Println("Node server listening on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept:", err)
			continue
		}

		go s.handleIncomingConnection(conn)
	}
}
