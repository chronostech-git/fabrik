package blockchain

import (
	"log"
	"time"

	"github.com/chronostech-git/fabrik/internal/blockchain/debug"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

type BlockCache []*Block

type ChainWriter interface {
	AddBlock(b *Block)          // Add a single block to the cache
	AddChain(blocks []*Block)   // Add multiple blocks to the cache
	FlushChainFromCache() error // Flush multiple blocks from the cache
}

type Chain struct {
	DB            storage.Database // NOTE: Using the interface means it can be leveldb OR memorydb.
	State         *state.ChainState
	BlockCache    BlockCache // is equal to []*Block
	Genesis       *Genesis
	Head          *Block // Current block
	ChainIterator storage.Iterator
	Validator     *Validator
	Printer       debug.Printer // Used for printing large quantities of data nicely (for debug purposes)

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
func NewWithGenesis(db storage.Database, coinbaseTx *Transaction, genesis *Genesis, gasLimit uint64) *Chain {
	c := New(db)
	c.Genesis = genesis
	c.State.AddBalance(genesis.Coinbase, genesis.InitialValue)

	if err := c.applyGenesis(coinbaseTx, gasLimit); err != nil {
		log.Panic(err)
	}

	return c
}

// If used, set the debug.Printer
func (c *Chain) SetDebugPrinter(printer debug.Printer) {
	c.Printer = printer
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

	b.Header.Height += 1
	c.Head = b

	return nil
}

// NOTE: This function is needed due to Block and Genesis technically being seperate structures.
// Yes, block fields and genesis fields may share similarities and are assumed to be the same type,
// this is why it is important to distinguish the difference between a genesis block, and any proceeding
// blocks being added to the chain.
func (c *Chain) applyGenesis(coinbaseTx *Transaction, gasLimit uint64) error {
	if c.Genesis == nil {
		return ErrGenesisMissing
	}

	c.Genesis.Txs = append(c.Genesis.Txs, coinbaseTx)

	genesisBlock := NewBlock(types.Empty32(), time.Now().Unix(), c.Genesis.Txs, 0, gasLimit)
	// if err := c.ApplyBlock(genesisBlock); err != nil {
	// 	return err
	// }

	c.Head = genesisBlock

	return nil
}

func (c *Chain) HasGenesis() bool {
	return c.Genesis != nil
}

func (c *Chain) GetBalance(addr types.Address) types.Amount {
	return c.State.GetBalance(addr)
}

// Returns true if the cache is empty
func (c *Chain) CacheEmpty() bool {
	return len(c.BlockCache) == 0
}

// Remove all blocks from cache
func (c *Chain) ClearCache() {
	clear(c.BlockCache)
}

func (c *Chain) writeBlock(b *Block) error {
	return b.Write(c.DataDir)
}

// Add a block to the cache
func (c *Chain) AddBlock(b *Block) {
	c.BlockCache = append(c.BlockCache, b)
}

// Add multiple blocks (chain of blocks) to the cache
func (c *Chain) AddChain(blocks []*Block) {
	for _, block := range blocks {
		c.AddBlock(block)
	}
}

// If blocks are in the cache, this function will commit them to the disk, and DB
// after it is already validated.
func (c *Chain) FlushChainFromCache() error {
	if c.CacheEmpty() {
		return ErrCacheEmpty
	}

	for _, block := range c.BlockCache {
		err := c.writeBlock(block)
		if err != nil {
			return ErrFailedWrite
		}

		blockBytes, err := rlp.Encode(block)
		if err != nil {
			return err
		}

		err = c.DB.Put(block.Hash.Bytes(), blockBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) Height() uint64 {
	return c.Head.Header.Height
}

// TODO This is where PrintPretty was residing. Create a "printer.go"
// which will print all chain state, chain data, and account data
