package blockchain

import (
	"fmt"
	"log"
	"time"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	MaxFutureBlockTime = 15 // We use this to calculate if a block is stale or not. If creationTime > creationTime + MaxFutureBlockTime -- it is considered stale data and will not be added to the chain.
)

type Chain struct {
	DB         storage.Database // NOTE: Using the interface means it can be leveldb OR memorydb.
	State      *state.ChainState
	BlockCache []*Block
	Genesis    *Genesis
	Head       *Block

	// For writing to disk. All block files will be stored in <datadir>.
	DataDir string
}

// Create a new empty chain WITHOUT genesis block.
func New(db storage.Database) *Chain {
	chain := &Chain{
		DB:    db,
		State: state.NewChainState(),
	}

	return chain
}

// Create a new chain WITH the genesis block.
func NewWithGenesis(db storage.Database, genesis *Genesis) *Chain {
	c := New(db)
	c.Genesis = genesis

	// For genesis state
	c.State.AddBalance(genesis.Coinbase, genesis.InitialValue)

	return c
}

// Set's the datadir to given path (param datadir).
func (c *Chain) SetDataDir(datadir string) {
	c.DataDir = datadir
}

// Updates the blockchain state (balances, accounts, transactions etc...)
// See state/chain_state.go for more information on how balances etc work
func (c *Chain) ApplyBlock(b *Block) error {
	txs := b.ToStateTxs()

	if err := c.State.ApplyTransactions(txs); err != nil {
		return err
	}

	c.Head = b

	return nil
}

// NOTE: This function is needed due to Block and Genesis technically being seperate structures.
// Yes, block fields and genesis fields may share similarities and are assumed to be the same type,
// this is why it is important to distinguish the difference between a genesis block, and any proceeding
// blocks being added to the chain.
func (c *Chain) ApplyGenesis(coinbaseTx *Transaction) error {
	if c.Genesis == nil {
		return ErrGenesisMissing
	}

	c.Genesis.Txs = append(c.Genesis.Txs, coinbaseTx)

	genesisBlock := NewBlock(types.Empty32(), time.Now().Unix(), c.Genesis.Txs, 0)
	if err := c.ApplyBlock(genesisBlock); err != nil {
		return err
	}

	c.Head = genesisBlock

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

// Flush (or push) contents of BlockCache and write it to disk / save to leveldb or memorydb.
// If cache is empty, it will panic.
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

// PrintPretty is used when --dump is called with cli/chain command.
// Example: cli/chain [...] --new --dump.
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
