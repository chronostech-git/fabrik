package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	DataDir string   `arg:"required" help:"data directory where all persistent data is stored."`
	Host    string   `arg:"--ipaddr" help:"host to bind the server" default:"127.0.0.1"`
	Port    string   `arg:"--port,required" help:"port to bind the server"`
	Peers   []string `arg:"--connect,separate" help:"peer addresses to connect to, e.g., 127.0.0.1:8000"`
}

// handleIncoming continuously reads messages from a peer connection
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
	if err := scanner.Err(); err != nil {
		log.Printf("[Peer %s] connection closed with error: %v\n", peer.ID, err)
	}
}

// startServer launches the TCP server for incoming peer connections
func startServer(addr string, mgr *p2p.PeerManager, disk *p2p.DiskStorage) {
	log.Printf("[FABNET] Server started on %s\n", addr)
	p2p.StartServer(addr, mgr, disk)
}

// connectToPeers dials a list of peer addresses and registers them in the PeerManager
func connectToPeers(peerAddrs []string, mgr *p2p.PeerManager, disk *p2p.DiskStorage) {
	for _, addr := range peerAddrs {
		peer, err := p2p.DialPeer(addr, mgr)
		if err != nil {
			log.Println("Failed to connect to peer", addr, err)
			continue
		}
		if err := disk.WritePeer(peer); err != nil {
			log.Println("Failed to save peer to disk:", err)
		}
		// Launch message listener
		go handleIncoming(peer)
	}
}

// registerExistingPeers writes all known peers to disk and starts listeners
func registerExistingPeers(peers map[string]*p2p.Peer, disk *p2p.DiskStorage) {
	for _, peer := range peers {
		if err := disk.WritePeer(peer); err != nil {
			log.Println("Failed to save peer to disk:", err)
			continue
		}
		go handleIncoming(peer)
	}
}

func main() {
	// Parse CLI arguments
	arg.MustParse(&args)

	disk := p2p.NewDiskStorage(args.DataDir)
	peermgr := p2p.NewPeerManager()
	addr := net.JoinHostPort(args.Host, args.Port)

	go startServer(addr, peermgr, disk)

	connectToPeers(args.Peers, peermgr, disk)
	registerExistingPeers(peermgr.Peers, disk)

	// Block forever (or implement a proper signal handler)
	select {}
}
