package main

import (
	"log"
	"net"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	Host string `arg:"--ipaddr"`
	Port string `arg:"required"`
}

func main() {
	arg.MustParse(&args)

	peermgr := p2p.NewPeerManager()

	addr := net.JoinHostPort(args.Host, args.Port)
	log.Printf("[FABNET] Server started ")

	p2p.StartServer(addr, peermgr)
}
