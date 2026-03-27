package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

var args struct {
	DataDir string `arg:"required"`
}

func main() {
	arg.MustParse(&args)

	if err := os.MkdirAll(args.DataDir, 0700); err != nil {
		log.Panic(err)
	}

	ks := keystore.NewFileStore(args.DataDir)

	wallet := blockchain.NewWallet(ks)
	account := wallet.CreateExternalAccount()

	fmt.Println("Address:", account.Address().Hex())
	fmt.Println("Key stored in keystore")
}
