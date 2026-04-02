package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	Connect string
}

func main() {
	arg.MustParse(&args)

	peermgr := p2p.NewPeerManager()

	if err := p2p.DialPeer(args.Connect, peermgr); err != nil {
		log.Panic(err)
	}

	for _, peer := range peermgr.Peers {
		go func(p *p2p.Peer) {
			scanner := bufio.NewScanner(p.Conn)
			for scanner.Scan() {
				msg, err := p2p.ParseMessage(scanner.Text())
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Printf("<%s> [%s]: %s\n", p.ID, msg.Type, msg.Data)
			}
		}(peer)
	}

	for _, peer := range peermgr.Peers {
		if err := peer.Send(&p2p.Message{Type: "PING", Data: ""}); err != nil {
			log.Println(err)
		}
	}

	select {}
}
