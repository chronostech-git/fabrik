package blockchain

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/types"
)

type BlockHeader struct {
	PrevHash  types.Hash
	Timestamp int64
	TxRoot    types.Hash
}

type Block struct {
	Header BlockHeader
	Txs    []*Transaction

	Hash types.Hash
}

func NewBlock(prevHash types.Hash, timestamp int64, txs []*Transaction) *Block {
	b := &Block{
		Header: BlockHeader{
			PrevHash:  prevHash,
			Timestamp: timestamp,
			TxRoot:    calcTxRoot(txs),
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
