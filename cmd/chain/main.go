package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/consensus"
	"github.com/chronostech-git/fabrik/internal/consensus/hawk"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/storage/leveldb"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir   string
	GasLimit  int
	Mechanism string // Hawk (PoW) or Falcon (PoI) -- "--mechanism <hawk-pow|falcon-poi>"
	UseMemory bool
}

// selectConsensusMechanism creates a new consensus engine base on the mechanism passed to it.
// Consensus options: "hawk-pow" | "falcon-poi"
func selectConsensusMechanism(mechanism string) consensus.Engine {
	switch mechanism {
	case "hawk-pow":
		return hawk.NewPoW()
	case "falcon-poi":
		return nil // TODO Implement Proof of Importance (PoI) consensus mechanism
	default:
		log.Println("Invalid consensus mechanism")
	}

	return nil
}

func loadKeyFromFile(datadir string) *crypto.Key {
	keystore := keystore.NewFileStore(datadir)

	key, err := keystore.GetKey()
	if err != nil {
		log.Panic(err)
	}

	return key
}

func createCoinbaseTx(datadir string, value types.Amount) *blockchain.Transaction {
	key := loadKeyFromFile(datadir)
	coinbaseTx := blockchain.NewTx(key.Address, types.ZeroAddress(), value, 0, nil)
	return coinbaseTx
}

func createGenesisBlock(
	currentTime int64,
	coinbaseTx *blockchain.Transaction,
) *blockchain.Genesis {
	return blockchain.NewGenesis(currentTime, coinbaseTx.Sender, coinbaseTx.Value)
}

func initChain(
	datadir string,
	gasLimit int,
	coinbaseValue types.Amount,
	useMemory bool,
) (*blockchain.Chain, *blockchain.Genesis, error) {
	var db storage.Database
	if useMemory {
		db = memory.New()
	}

	ldbFile := filepath.Join(datadir, "ldb-manifest")
	db, err := leveldb.New(ldbFile)
	if err != nil {
		return nil, nil, err
	}

	coinbaseTx := createCoinbaseTx(datadir, coinbaseValue)
	genesisBlock := createGenesisBlock(time.Now().Unix(), coinbaseTx)

	chain := blockchain.NewWithGenesis(
		db,
		coinbaseTx,
		genesisBlock,
		uint64(gasLimit),
	)

	return chain, genesisBlock, nil
}

func main() {
	arg.MustParse(&args)

	chain, genesis, err := initChain(
		args.DataDir,
		args.GasLimit,
		types.ZeroAmount(),
		args.UseMemory,
	)
	if err != nil {
		log.Panic(err)
	}

	chain.SetDataDir(args.DataDir)

	if err = genesis.Write(args.DataDir); err != nil {
		log.Panic(err)
	}

	chosenMechanism := selectConsensusMechanism(args.Mechanism)
	chain.SetConsensusMechanism(chosenMechanism)

	prevBlock, err := genesis.ToBlock()
	if err != nil {
		log.Panic(err)
	}

	prevBlockView := prevBlock.ToConsensusBlockView()

	newTestBlock := blockchain.NewBlock(
		prevBlock.Hash,
		time.Now().Unix(),
		chain.Head.Txs,
		chain.Height(),
		uint64(args.GasLimit),
	)
	chain.AddBlock(newTestBlock)

	if err := chain.FlushCache(); err != nil {
		log.Panic(err)
	}

	switch args.Mechanism {
	case "hawk-pow":
		difficultyTarget := hawk.CalcPoWDifficulty(
			prevBlockView.DifficultyTarget,
			prevBlockView.Header.Timestamp,
			chain.ToConsensusChainView().Head.Header.Timestamp,
			chain.Head.Header.Timestamp+blockchain.MaxFutureBlockTime,
		)
		chain.Engine.RunPoW(
			difficultyTarget,
			newTestBlock.ToConsensusBlockView(),
		)
	case "falcon-poi":
		fmt.Println("falcon-poi not implemented yet")
	}
}
