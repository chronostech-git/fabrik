package main

import (
	"crypto/sha256"
	"log"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	blockprinter "github.com/chronostech-git/fabrik/internal/blockchain/debug/block_printer"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/storage/leveldb"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir         string `arg:"required"`
	GenisisGasLimit int    `arg:"--gaslimit"`
	New             bool
	Debug           bool
}

func signAndVerifyCoinbaseTx(coinbaseTx *blockchain.Transaction, key *crypto.Key) (*crypto.Signature, bool, error) {
	data, err := rlp.Encode(coinbaseTx)
	if err != nil {
		return nil, false, err
	}

	dataHash := sha256.Sum256(data)

	sig, err := key.Sign(dataHash)
	if err != nil {
		return nil, false, err
	}

	validSig := key.Verify(dataHash, sig)

	if validSig {
		coinbaseTx.Signature = sig
	}

	return sig, validSig, nil
}

func commitGenesisBlock(chain *blockchain.Chain, genesisBlock *blockchain.Block) error {
	chain.AddBlock(genesisBlock)
	if err := chain.FlushChainFromCache(); err != nil {
		return err
	}
	log.Printf("Genesis block commited to chain with hash %s", genesisBlock.Hash.String())
	return nil
}

func main() {
	arg.MustParse(&args)

	keystore := keystore.NewFileStore(args.DataDir)

	key, err := keystore.GetKey()
	if err != nil {
		log.Panic(err)
	}

	// Genesis funding transaction
	coinbase := blockchain.NewTx(
		types.ZeroAddress(),
		key.Address,
		types.NewAmount(1_000_000),
		1,
		nil,
	)

	sig, valid, err := signAndVerifyCoinbaseTx(coinbase, key)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Coinbase transaction signed and verified\n\tsigner: %s\n\tsig: %s\n\thash: %s\n\tvalid: %t\n",
		key.Address.String(), sig.Hex(), coinbase.Hash.String(), valid)

	genesis := blockchain.NewGenesis(
		time.Now().Unix(),
		coinbase.Sender,
		types.NewAmount(1),
	)

	ldbFile := filepath.Join(args.DataDir, "ldb")
	ldb, err := leveldb.New(ldbFile)
	if err != nil {
		log.Panic(err)
	}

	chain := blockchain.NewWithGenesis(
		ldb,
		coinbase,
		genesis,
		uint64(args.GenisisGasLimit),
	)

	chain.SetDataDir(args.DataDir)

	log.Printf("Genesis block created at %d unix time with hash %s",
		chain.Genesis.CreationTime, chain.Genesis.GenesisHash.String())

	genesisToBlock, err := genesis.ToBlock()
	if err != nil {
		log.Panic(err)
	}

	if err := commitGenesisBlock(chain, genesisToBlock); err != nil {
		log.Panic(err)
	}

	if args.Debug {
		blockPrinter := blockprinter.New()
		blockPrinter.SetBlock(genesisToBlock)
		chain.SetDebugPrinter(blockPrinter)
		chain.Printer.PrintData()
	}
}
