package consensus

import (
	"crypto/sha256"
	"log"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

// HeaderView is a consensus compatible version of
// a BlockHeader.
type HeaderView struct {
	PrevHash  types.Hash
	Timestamp int64
	TxRoot    types.Hash
	Height    uint64
}

// BlockView is a consensus compatible version of a
// Block.
type BlockView struct {
	Header *HeaderView
	Hash   types.Hash
	Txs    []state.Tx

	// Added for hawk PoW consensus mechanism
	Nonce            uint64
	DifficultyTarget uint64
	Miner            types.Address
	MinerReward      types.Amount
	Mined            bool
}

// CalcHawkHash is used as a differing hash calculation
// from (*block).Hash. Notice, the fields are different
// between a blockchain.Block and consensus.BlockView. This
// is because PoW will add fields, and have overall different fields
// than a base blockchain.Block
func (bv *BlockView) CalcHawkHash() types.Hash {
	enc, err := rlp.Encode(bv)
	if err != nil {
		log.Panic(err)
	}

	return sha256.Sum256(enc)
}
