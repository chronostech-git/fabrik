package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	ServerAddress string `arg:"--fabnet"`
}

func main() {
	arg.MustParse(&args)

	peermgr := p2p.NewPeerManager()

	if err := p2p.DialPeer(args.ServerAddress, peermgr); err != nil {
		log.Panic(err)
	}

	log.Printf("[NODE] Dialing peer %s...", args.ServerAddress)

	select {}
}
