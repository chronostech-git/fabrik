package blockchain

import (
	"fmt"
	"log"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	MaxFutureBlockTime = 15
)

type Chain struct {
	DB         storage.Database
	State      *state.ChainState
	BlockCache []*Block
	Genesis    *Genesis
	Head       *Block

	// For disk
	DataDir string
}

func New(db storage.Database) *Chain {
	chain := &Chain{
		DB:    db,
		State: state.NewChainState(),
	}

	return chain
}

func NewWithGenesis(db storage.Database, genesis *Genesis) *Chain {
	c := New(db)
	c.Genesis = genesis

	// For genesis state
	c.State.AddBalance(genesis.Coinbase, genesis.InitialValue)

	return c
}

func (c *Chain) SetDataDir(datadir string) {
	c.DataDir = datadir
}

func (c *Chain) ApplyBlock(b *Block) error {
	// TODO: uncomment this when ready
	// validBlock, err := c.Validator.ValidateBlock(b)
	// if err != nil {
	// 	return err
	// }

	txs := b.ToStateTxs()

	if err := c.State.ApplyTransactions(txs); err != nil {
		return err
	}

	c.Head = b

	return nil
}

func (c *Chain) GetBalance(addr types.Address) types.Amount {
	return c.State.GetBalance(addr)
}

func (c *Chain) emptyCache() bool {
	return len(c.BlockCache) == 0
}

func (c *Chain) writeBlock(b *Block) error {
	return b.Write(c.DataDir)
}

func (c *Chain) AddBlockToCache(b *Block) {
	c.BlockCache = append(c.BlockCache, b)
}

func (c *Chain) FlushCacheToDisk() error {
	if c.emptyCache() {
		return ErrCacheEmpty
	}

	for i, block := range c.BlockCache {
		err := c.ApplyBlock(block)
		if err != nil {
			return err
		}

		err = c.writeBlock(block)
		if err != nil {
			return err
		}

		blockData, err := rlp.Encode(block)
		if err != nil {
			return err
		}
		c.DB.Put(block.Hash.Bytes(), blockData)

		log.Println("Block #%05d flushed to disk\n", i)
	}

	fmt.Println("Finished.")

	return nil
}

func (c *Chain) PrintPretty() {
	fmt.Println("Genesis Data")
	fmt.Println("\thash:", c.Genesis.GenesisHash.String())
	fmt.Println("\tvalue:", c.Genesis.InitialValue.String())
	fmt.Println()
	fmt.Println("Current Block Data")
	fmt.Println("\thash:", c.Head.Hash.String())
	fmt.Println("\ttime:", c.Head.Header.Timestamp)
	fmt.Println("\ttxroot:", c.Head.Header.TxRoot.String())
	fmt.Println("\theight:", c.Head.Header.Height)
	fmt.Println()
	fmt.Println("State Balance Data")
	numAcc := 0
	for addr, bal := range c.State.Balances() {
		fmt.Printf("\tAccount #%d\n", numAcc+1)
		fmt.Println("\t\taddr:", addr.String())
		fmt.Println("\t\tbalance:", bal.String())
		fmt.Println()
		numAcc++
	}

}
