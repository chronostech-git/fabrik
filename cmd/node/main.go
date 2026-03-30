package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/storage/leveldb"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir   string
	UseMemory bool
	Debug     bool
	Dump      bool
}

func main() {
	arg.MustParse(&args)

	var db storage.Database
	if args.UseMemory {
		db = memory.New()
	}

	leveldbPath := filepath.Join(args.DataDir, "manifest")

	db, err := leveldb.New(leveldbPath)
	if err != nil {
		log.Panic(err)
	}

	genesis, err := blockchain.LoadGenesis(args.DataDir)
	if err != nil {
		log.Panic(err)
	}

	chain := blockchain.NewWithGenesis(db, genesis)
	chain.SetDataDir(args.DataDir)

	var txs []*blockchain.Transaction

	tx1 := blockchain.NewTx(genesis.Coinbase, types.ZeroAddress(), types.NewAmount(10000), 0, nil)

	txs = append(txs, tx1)

	newBlock := blockchain.NewBlock(
		genesis.GenesisHash,
		time.Now().Unix(),
		txs,
		1,
	)

	chain.AddBlockToCache(newBlock)

	if err := chain.FlushCacheToDisk(); err != nil {
		log.Panic(err)
	}

	if args.Dump {
		chain.PrintPretty()
	}

}
