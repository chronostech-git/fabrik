package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/storage/leveldb"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir string `arg:"required"`
	New     bool   // If true, chain will automatically generate a genesis block under the supplied Data Directory
	Memory  bool   // If true, it will use memory.New() instead of leveldb.New() (not persistent v. persistent storage)
	Dump    bool   // If true, it will call chain.PrintPretty()
}

func createNewChainWithGenesis(datadir string, coinbaseKey *crypto.Key, useMemory bool) *blockchain.Chain {
	var chain *blockchain.Chain
	var db storage.Database

	currentTime := time.Now().Unix()
	initialCirculatingValue := types.NewAmount(1_000_000_000_000_000)
	genesisBlock := blockchain.NewGenesis(currentTime, coinbaseKey.Address, initialCirculatingValue)

	if useMemory {
		db = memory.New()
		chain = blockchain.NewWithGenesis(db, genesisBlock)
	} else {
		db, err := leveldb.New(datadir)
		if err != nil {
			log.Panic(err)
		}
		chain = blockchain.NewWithGenesis(db, genesisBlock)
	}

	return chain
}

func createCoinbaseTransaction(coinbaseKey *crypto.Key) *blockchain.Transaction {
	coinbaseTx := blockchain.NewTx(types.ZeroAddress(), coinbaseKey.Address,
		types.NewAmount(1_000_000_000_000_000), 0, nil)

	// sign transaction manually (without wallet) using the Coinbase Transaction key.
	// A coinbase requires a wallet to be created, therfore no need to pass "keystore"
	// as an argument here...it implies a wallet has already been created if coinbaseKey != nil.
	// So, we can just pass the coinbase key directly.
	sig, err := coinbaseKey.Sign(coinbaseTx.Hash)
	if err != nil {
		log.Panic(err)
	}

	validSig := coinbaseKey.Verify(coinbaseTx.Hash, sig)

	log.Printf("Coinbase transaction created, signed, and verified:\n\t hash=%s\n\t sig=%s\n\t valid=%t\n\n",
		coinbaseTx.Hash.String(), sig.Hex(), validSig)

	coinbaseTx.Signature = sig

	return coinbaseTx
}

func main() {
	arg.MustParse(&args)

	if !args.New {
		fmt.Printf("ERROR: Chain already created with genesis block under %s/genesis directory", args.DataDir)
		return
	}

	keystore := keystore.NewFileStore(args.DataDir)

	coinbaseKey, err := keystore.GetKey()
	if err != nil {
		log.Panic(err)
	}

	chain := createNewChainWithGenesis(args.DataDir, coinbaseKey, args.Memory)

	coinbase := createCoinbaseTransaction(coinbaseKey)
	if err := chain.ApplyGenesis(coinbase); err != nil {
		log.Panic(err)
	}

	if args.Dump {
		chain.PrintPretty()
	}

	fmt.Println("Finished.")
}
