package blockchain

import (
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

type Chain struct {
	DB           storage.Database
	State        *state.ChainState
	CurrentBlock *Block
	Genesis      *Genesis
	Head         *Block
}

func New(db storage.Database) *Chain {
	return &Chain{
		DB:    db,
		State: state.NewChainState(),
	}
}

func NewWithGenesis(db storage.Database, genesis *Genesis) *Chain {
	c := New(db)
	c.Genesis = genesis

	// For genesis state
	c.State.AddBalance(genesis.Coinbase, genesis.InitialValue)

	return c
}

func (c *Chain) ApplyBlock(b *Block) error {
	// This is extremely basic validation
	// TODO: Implement block validators
	if c.Head != nil {
		if b.Header.PrevHash != c.Head.Hash {
			return ErrInvalidBlock
		}
	}

	txs := b.ToStateTxs()

	if err := c.State.ApplyTransactions(txs); err != nil {
		return err
	}

	c.Head = b

	return nil
}

func (c *Chain) GetBalance(addr types.Address) types.Amount {
	return c.State.GetBalance(addr)
}
