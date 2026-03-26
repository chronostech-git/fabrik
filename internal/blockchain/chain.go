package blockchain

import (
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
)

type Chain struct {
	DB storage.Database

	State *state.ChainState

	CurrentBlock *Block
	Genesis      *Genesis
}

func New(db storage.Database) *Chain {
	return &Chain{
		DB:    db,
		State: state.NewChainState(db),
	}
}

func NewWithGenesis(db storage.Database, genesis *Genesis) *Chain {
	c := New(db)
	c.Genesis = genesis
	return c
}
