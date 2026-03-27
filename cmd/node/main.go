package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/storage/leveldb"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
)

var args struct {
	DataDir   string
	UseMemory bool
	// Many other args are needed, this is just the beginning...Mwahahahah
}

func main() {
	arg.MustParse(&args)

	genesis, err := blockchain.LoadGenesis(args.DataDir)
	if err != nil {
		log.Panic(err)
	}

	var db storage.Database
	if args.UseMemory {
		db = memory.New()
	} else {
		path := filepath.Join(args.DataDir, "ondisk")
		db, err = leveldb.New(path)
		if err != nil {
			log.Panic(err)
		}
	}

	chain := blockchain.NewWithGenesis(db, genesis)

	fmt.Printf("Genesis Coinbase: %x\n", genesis.Coinbase)

	for k, v := range chain.State.Balances() {
		fmt.Printf("Stored Key: %x Balance: %v\n", k, v)
	}

	balance := chain.State.GetBalance(genesis.Coinbase)
	fmt.Println(balance.String())
}
