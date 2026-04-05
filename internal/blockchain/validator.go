package blockchain

import (
	"log"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/types"
	"github.com/google/uuid"
)

type Validator struct {
	identity uuid.UUID
	chain    *Chain // Every validator has a copy of the blockchain
}

func NewValidator(chainCopy *Chain) *Validator {
	identityAddr, err := uuid.NewV6()
	if err != nil {
		log.Panic(err)
	}

	return &Validator{
		identity: identityAddr,
		chain:    chainCopy,
	}
}

func (v *Validator) IdentityString() string {
	return v.identity.String()
}

func (v *Validator) GetBlockByHeight(height uint64) *Block {
	iter := v.chain.ChainIterator

	var block Block

	for iter.Next() {
		blockBytes := iter.Value()
		err := rlp.Decode(blockBytes, &block)

		if err != nil {
			log.Panic(err)
			break
		}

		if block.Header.Height == height {
			return &block
		}
	}

	return nil
}

func (v *Validator) GetBlockByHash(hash types.Hash) *Block {
	blockBytes, err := v.chain.DB.Get(hash.Bytes())
	if err != nil {
		log.Panic(err)
	}

	var block Block
	err = rlp.Decode(blockBytes, &block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (v *Validator) ValidateBlockHeader(h *BlockHeader) (bool, error) {
	if h.PrevHash.IsZero() {
		return false, ErrInvalidPrevHash
	}
	if h.TxRoot.IsZero() {
		return false, ErrNoTransactions
	}
	return true, nil
}

func (v *Validator) ValidateBlockBody(b *Block) (bool, error) {
	validHeader, err := v.ValidateBlockHeader(&b.Header)
	if err != nil || !validHeader {
		return false, err
	}
	if b.Hash.IsZero() {
		return false, ErrInvalidBlock
	}

	// TODO Check if StateRoot.IsZero() here when implemented...
	if !b.HasTxs() {
		return false, ErrNoTransactions
	}
	if b.IsStale() {
		return false, ErrMaxBlockTimeReached // marked as "stale"
	}
	if b.Size() >= MaxBlockSize {
		return false, ErrBlockSizeReached
	}

	return true, nil
}

func (v *Validator) IsValidBlock(b *Block) (bool, error) {
	validHeader, err := v.ValidateBlockHeader(&b.Header)
	if err != nil {
		return false, err
	}
	if !validHeader {
		return false, err
	}
	validBody, err := v.ValidateBlockBody(b)
	if err != nil {
		return false, err
	}
	if !validBody {
		return false, err
	}

	return true, nil
}

func (v *Validator) CommitBlock(b *Block) error {
	// Add block to cache
	v.chain.AddBlock(b)

	// Apply block to state
	err := v.chain.ApplyBlock(b)
	if err != nil {
		return err
	}

	// Immediately flush the cache to the disk
	err = v.chain.FlushChainFromCache()
	if err != nil {
		return err
	}

	log.Printf("Block #%08d validated at %d unix time and commited with hash %s",
		b.Header.Height, b.Header.Timestamp, b.Hash.String())

	return nil
}
