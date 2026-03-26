package state

import "github.com/chronostech-git/fabrik/internal/storage"

type ChainState struct {
	db storage.Database
}

func NewChainState(db storage.Database) *ChainState {
	return &ChainState{db: db}
}
