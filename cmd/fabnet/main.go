package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	Host string `arg:"--ipaddr"`
	Port string `arg:"required"`
}

func buildServerAddress(host, port string) string {
	return host + ":" + port
}

func main() {
	arg.MustParse(&args)

	peermgr := p2p.NewPeerManager()

	addr := buildServerAddress(args.Host, args.Port)

	log.Printf("[FABNET] Server starting... You may now connect using cli/node!")
	p2p.StartServer(addr, peermgr)
}
