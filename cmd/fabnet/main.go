package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	Host  string   `arg:"--ipaddr" help:"host to bind the server" default:"127.0.0.1"`
	Port  string   `arg:"--port,required" help:"port to bind the server"`
	Peers []string `arg:"--connect,separate" help:"peer addresses to connect to, e.g., 127.0.0.1:8000"`
}

func handleIncoming(peer *p2p.Peer) {
	scanner := bufio.NewScanner(peer.Conn)
	for scanner.Scan() {
		line := scanner.Text()
		msg, err := p2p.ParseMessage(line)
		if err != nil {
			log.Println("Parse error:", err)
			continue
		}
		fmt.Printf("[From %s] %s: %s\n", peer.ID, msg.Type, msg.Data)
	}
}

func main() {
	arg.MustParse(&args)

	peermgr := p2p.NewPeerManager()
	addr := net.JoinHostPort(args.Host, args.Port)

	// Start server
	go func() {
		log.Printf("[FABNET] Server started on %s\n", addr)
		p2p.StartServer(addr, peermgr)
	}()

	// Connect to any peers specified
	for _, peerAddr := range args.Peers {
		if err := p2p.DialPeer(peerAddr, peermgr); err != nil {
			log.Println("Failed to connect to peer", peerAddr, err)
		}
	}

	// Launch listeners for all connected peers
	for _, peer := range peermgr.Peers {
		go handleIncoming(peer)
	}

	// Interactive CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		// Send to all connected peers
		for _, peer := range peermgr.Peers {
			err := peer.Send(&p2p.Message{Type: text, Data: ""})
			if err != nil {
				log.Println("Send error:", err)
			}
		}
	}
}
