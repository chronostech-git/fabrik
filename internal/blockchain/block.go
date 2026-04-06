package blockchain

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/blockchain/consensus"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

type BlockHeader struct {
	PrevHash  types.Hash   // Previous block hash
	Timestamp int64        // Time of creation
	TxRoot    types.Hash   // All transactions present in block are then hashed and stored in the TxRoot
	Height    uint64       // Block number
	GasLimit  uint64       // Max amount of gas to be used during block creation / application
	GasUsed   uint64       // Total amount of gas burned for specific block
	BaseFee   types.Amount // Base fee per gas in FAB
}

type Block struct {
	Header BlockHeader    // Header data
	Txs    []*Transaction // List of transactions packaged in block

	// NOTE: State root is currently not implemented
	// The "state root" is designed to make sure all nodes agree on the exact
	// state of Fabrik at any given block.
	StateRoot types.Hash

	Hash types.Hash
}

// NewBlock create a new blockchain block
func NewBlock(prevHash types.Hash, timestamp int64, txs []*Transaction, height uint64, gasLimit uint64) *Block {
	b := &Block{
		Header: BlockHeader{
			PrevHash:  prevHash,
			Timestamp: timestamp,
			TxRoot:    calcTxRoot(txs),
			Height:    height,
			GasLimit:  gasLimit,
		},
		Txs: txs,
	}

	b.Hash = b.computeHash()
	return b
}

// computeHash calculate the hash of a block.
// Used in NewBlock
func (b *Block) computeHash() types.Hash {
	enc, _ := rlp.Encode(b.Header)
	return sha256.Sum256(enc)
}

// calcTxRoot calculate the root hash of the transactions list of a given block.
// Used in NewBlock() function.
func calcTxRoot(txs []*Transaction) types.Hash {
	enc, _ := rlp.Encode(txs)
	return sha256.Sum256(enc)
}

func (b *Block) Transactions() []state.TxView {
	out := make([]state.TxView, len(b.Txs))
	for i, tx := range b.Txs {
		out[i] = tx
	}
	return out
}

// ToStateTxs converts all transactions within a block
// to a state-compatible TxView.
func (b *Block) ToStateTxs() []state.Tx {
	out := make([]state.Tx, len(b.Txs))

	for i, tx := range b.Txs {
		out[i] = state.Tx{
			From:  tx.Sender,
			To:    tx.Receiver,
			Value: tx.Value,
			Data:  []byte{},
		}
	}
	return out
}

// Write write a block to disk located in
// <datadir>/blocks/<block-filename>.dat file
func (b *Block) Write(datadir string) error {
	data, err := rlp.Encode(b)
	if err != nil {
		return err
	}

	dir := filepath.Join(datadir, "blocks")
	fmt.Println(dir)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path := filepath.Join(dir, BlockFilename(b))
	return os.WriteFile(path, data, 0600)
}

// CalcGasRemaining calculates the total gas remaining
// after block transaction execution.
func (b *Block) CalcGasRemaining() uint64 {
	gasRemaining := b.Header.GasLimit - b.Header.GasUsed
	if gasRemaining <= 0 {
		log.Panicf("Out of gas for block #%08d", b.Header.Height)
	}
	return gasRemaining
}

func (b *Block) HasTxs() bool {
	return len(b.Txs) != 0
}

func (b *Block) IsStale() bool {
	return b.Header.Timestamp >= b.Header.Timestamp+MaxFutureBlockTime
}

// Size returns the length of an encoded block at time of call.
// The max block size is 2_000_000 bytes (2 MB).
func (b *Block) Size() int {
	enc, err := rlp.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return len(enc)
}

// Consensus functions
// consensus uses BlockView and ChainView as consensus-compatible
// versions of blockchain.Block and blockchain.Chain.

// ToConsensusBlockView converts a traditional blockchain.Block into
// a consensus-compatible BlockView.
func (b *Block) ToConsensusBlockView() *consensus.BlockView {
	header := &consensus.HeaderView{
		PrevHash:  b.Header.PrevHash,
		Timestamp: b.Header.Timestamp,
		TxRoot:    b.Header.TxRoot,
		Height:    b.Header.Height,
	}

	return &consensus.BlockView{
		Header:           header,
		Hash:             b.Hash,
		Txs:              b.ToStateTxs(),
		Nonce:            0,
		DifficultyTarget: 0,
		Miner:            types.ZeroAddress(),
		MinerReward:      types.ZeroAmount(),
		Mined:            false,
	}
}
