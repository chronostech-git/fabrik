package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir string `arg:"required"`
	Extra   string `arg:"required"`
}

func main() {
	arg.MustParse(&args)

	ks := keystore.NewFileStore(args.DataDir)

	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	tx := blockchain.NewTx(
		key.Address,
		types.Address{},
		types.NewAmount(1_000_000),
		0,
		[]byte(args.Extra),
	)

	sig, err := key.Sign(tx.Hash)
	if err != nil {
		log.Panic(err)
	}
	tx.Signature = sig

	genesis := blockchain.NewGenesis(
		time.Now().Unix(),
		key.Address,
		tx.Value,
	)

	genesisDir := filepath.Join(args.DataDir, "genesis")

	if err := genesis.Write(genesisDir); err != nil {
		log.Panic(err)
	}

	log.Println("Genesis written to:", filepath.Join(genesisDir, "genesis.dat"))
}
