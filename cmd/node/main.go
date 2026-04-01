package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/p2p"
)

var args struct {
	Boot bool
	Port string
	Peer bool
}

func startNetworkServer() {
	server := p2p.NewServer()
	server.Start(":" + args.Port)
}

func connectToNetworkAsPeer() {
	conn, err := net.Dial("tcp", "localhost:7000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type messages:")

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			log.Println("Server:", scanner.Text())
		}
	}()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		text := input.Text()
		log.Println(conn.LocalAddr().String()+":", text)
	}
}

func main() {
	arg.MustParse(&args)

	if args.Boot {
		startNetworkServer()
	} else if args.Peer {
		connectToNetworkAsPeer()
	}

	log.Println("Network has stopped...")
}
