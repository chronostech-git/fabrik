package blockchain

import "github.com/chronostech-git/fabrik/internal/types"

type ChainReader interface {
	GetBlockByHeight(height uint64) *Block // Get a block by the number (incremented +1 every block added to chain)
	GetBlockByHash(hash types.Hash) *Block // Get a block by it's associated hash value
}

type ConsensusEngine interface {
	ChainReader

	// Validation functions
	ValidateBlockHeader(h *BlockHeader) (bool, error) // Validate only the block header
	ValidateBlockBody(b *Block) (bool, error)         // Validate everything including the block header

	IsValidBlock(b *Block) bool // Check if a block is valid by hash

	// Used once validation, and packaging is done. This writes
	// already-validated blocks into the state db / leveldb / disk
	CommitBlock(b *Block) error
}
