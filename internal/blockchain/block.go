package blockchain

import (
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

type BlockHeader struct {
	PrevHash  types.Hash
	Timestamp int64
	TxRoot    types.Hash
	Height    uint64
}

type Block struct {
	Header BlockHeader
	Txs    []*Transaction

	StateRoot types.Hash
	Hash      types.Hash
}

// Create a new block
func NewBlock(prevHash types.Hash, timestamp int64, txs []*Transaction, height uint64) *Block {
	b := &Block{
		Header: BlockHeader{
			PrevHash:  prevHash,
			Timestamp: timestamp,
			TxRoot:    calcTxRoot(txs),
			Height:    height,
		},
		Txs: txs,
	}

	b.Hash = b.computeHash()
	return b
}

// Compute, or calculate the blocks hash.
// Used in NewBlock() function
func (b *Block) computeHash() types.Hash {
	enc, _ := rlp.Encode(b.Header)
	return sha256.Sum256(enc)
}

// Calculate the root hash of the transactions list of a given block.
// Used in NewBlock() function.
func calcTxRoot(txs []*Transaction) types.Hash {
	enc, _ := rlp.Encode(txs)
	return sha256.Sum256(enc)
}

// BlockView interface function
func (b *Block) Transactions() []state.TxView {
	out := make([]state.TxView, len(b.Txs))
	for i, tx := range b.Txs {
		out[i] = tx
	}
	return out
}

// IMPORTANT:
// This function converts the transactions list of a given block
// into state transactions. It is here so that we can avoid circular imports
// when communicating between state/ and the blockchain/ folder.
// See state/transaction_view.go and state/chain_state.go to see how this works.
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

// Write the block to disk using a special filename consisting of <block hash>-<block height>.dat.
// NOTE: Blocks are also written to leveldb.
func (b *Block) Write(datadir string) error {
	data, err := rlp.Encode(b)
	if err != nil {
		return err
	}

	// format: ./<datadir>/blocks/<block-hash>-<block-height>.dat
	// example: ./data/blocks/00000000000000000000000000000000-00000000.dat
	if err := os.MkdirAll(datadir+"/blocks/", 0700); err != nil {
		return err
	}

	path := filepath.Join(datadir+"/blocks/", BlockFilename(b))
	return os.WriteFile(path, data, 0600)
}
