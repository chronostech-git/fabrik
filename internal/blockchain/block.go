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

func (b *Block) computeHash() types.Hash {
	enc, _ := rlp.Encode(b.Header)
	return sha256.Sum256(enc)
}

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
