package blockchain

import (
	"log"
	"time"

	"github.com/chronostech-git/fabrik/internal/blockchain/consensus"
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
	Engine        consensus.Engine
	Printer       debug.Printer // Used for printing large quantities of data nicely (for debug purposes)

	// For writing to disk. All block files will be stored in <datadir>.
	DataDir string
}

// New creates a new empty chain excluding the genesis block
func New(db storage.Database) *Chain {
	chain := &Chain{
		DB:    db,
		State: state.NewChainState(),
	}

	return chain
}

// NewWithGenesis creates an empty chain using New and sets the chains
// genesis block to the provided genesis when calling NewWithGenesis.
func NewWithGenesis(db storage.Database, coinbaseTx *Transaction, genesis *Genesis, gasLimit uint64) *Chain {
	c := New(db)
	c.Genesis = genesis
	c.State.AddBalance(genesis.Coinbase, genesis.InitialValue)

	if err := c.applyGenesis(coinbaseTx, gasLimit); err != nil {
		log.Panic(err)
	}

	return c
}

// SetConsensusMechanism is used in cmd/chain/main.go.
// When using cli/chain.exe command line tool, the user can
// choose between the two provided consensus mechanisms:
// "hawk-pow" 	- Proof Of Work consensus
// "falcon-poi" - Proof Of Importance consensus
func (c *Chain) SetConsensusMechanism(e consensus.Engine) {
	c.Engine = e
}

// SetDataDir is used in cmd/chain/main.go to set the chains
// data directory to the provided folder passed in cli/chain.exe.
func (c *Chain) SetDataDir(datadir string) {
	c.DataDir = datadir
}

// ApplyBlock applies the block to the state.
// Also increments the block height and sets the chain head
// to the block after all transactions within the block have been applied
// to the state.
func (c *Chain) ApplyBlock(b *Block) error {
	txs := b.ToStateTxs()

	if err := c.State.ApplyTransactions(txs); err != nil {
		return err
	}

	b.Header.Height += 1
	c.Head = b

	return nil
}

// applyGenesis does the same as ApplyBlock, however it is private
// due to only being called within the blockchain module.
// Note: Block and Genesis are not inherently the same. A
// Genesis is a Block...A Block is not a Genesis.
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

// HasGenesis returns true if the genesis is present in the chain.
func (c *Chain) HasGenesis() bool {
	return c.Genesis != nil
}

// CacheEmpty returns true if the cache has no blocks.
func (c *Chain) CacheEmpty() bool {
	return len(c.BlockCache) == 0
}

// ClearCache erases all blocks that may, or may not be in the cache.
func (c *Chain) ClearCache() {
	clear(c.BlockCache)
}

// writeBlock is a helper function used within FlushCache function.
// Writes provided block to the disk (.dat file) under the blocks/ folder
// in <datadir>.
func (c *Chain) writeBlock(b *Block) error {
	return b.Write(c.DataDir)
}

// AddBlock adds a block to the cache
func (c *Chain) AddBlock(b *Block) {
	c.BlockCache = append(c.BlockCache, b)
}

// FlushCache does the following-
// 1. Checks if the cache is currently empty
// 2. Iterates over all blocks within the block cache
// 3. Writes each iterated block to a .dat file
// 4. Encodes iterated block and writes it to leveldb
func (c *Chain) FlushCache() error {
	if c.CacheEmpty() {
		return ErrCacheEmpty
	}

	for _, block := range c.BlockCache {
		// Write block to disk
		err := c.writeBlock(block)
		if err != nil {
			return ErrFailedWrite
		}

		// Encode block
		blockBytes, err := rlp.Encode(block)
		if err != nil {
			return err
		}

		// Commit block to leveldb or memorydb
		err = c.DB.Put(block.Hash.Bytes(), blockBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Height returns the current height of the chain.
func (c *Chain) Height() uint64 {
	return c.Head.Header.Height
}

// ToConsensusChainView converts the chain into a consensus-compatible
// chain.
func (c *Chain) ToConsensusChainView() *consensus.ChainView {
	return &consensus.ChainView{
		DB:    c.DB,
		State: c.State,
		Head:  c.Head.ToConsensusBlockView(),
	}
}
