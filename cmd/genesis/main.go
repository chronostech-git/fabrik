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

	// ✅ load key from keystore
	ks := keystore.NewFileStore(args.DataDir)

	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	// ✅ build transaction properly
	tx := blockchain.NewTx(
		key.Address,
		types.Address{}, // zero = mint
		types.NewAmount(1_000_000),
		0,
		[]byte(args.Extra),
	)

	// ✅ sign tx
	sig, err := key.Sign(tx.Hash)
	if err != nil {
		log.Panic(err)
	}
	tx.Signature = sig

	// ✅ create genesis (PURE DATA)
	genesis := blockchain.NewGenesis(
		time.Now().Unix(),
		key.Address,
		tx.Value,
	)

	genesisDir := filepath.Join(args.DataDir, "genesis")

	if err := genesis.Write(genesisDir); err != nil {
		log.Panic(err)
	}
}
