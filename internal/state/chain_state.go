package state

import (
	"errors"

	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Tx struct {
	From  types.Address
	To    types.Address
	Value types.Amount
}

type ChainState struct {
	balances map[types.Address]types.Amount
}

func NewChainState() *ChainState {
	return &ChainState{
		balances: make(map[types.Address]types.Amount),
	}
}

func (cs *ChainState) GetBalance(addr types.Address) types.Amount {
	bal, ok := cs.balances[addr]
	if !ok {
		return types.NewAmount(0)
	}
	return bal
}

func (cs *ChainState) SetBalance(addr types.Address, amount types.Amount) {
	cs.balances[addr] = amount
}

func (cs *ChainState) AddBalance(addr types.Address, amount types.Amount) {
	current := cs.GetBalance(addr)
	cs.balances[addr] = current.Add(amount)
}

func (cs *ChainState) SubtractBalance(addr types.Address, amount types.Amount) error {
	current := cs.GetBalance(addr)

	if current.LessThan(amount) {
		return ErrInsufficientBalance
	}

	res, err := current.Sub(amount)
	if err != nil {
		return err
	}

	cs.balances[addr] = res

	return nil
}

func (cs *ChainState) ApplyTx(tx Tx) error {
	if tx.From.IsZero() {
		cs.AddBalance(tx.To, tx.Value)
		return nil
	}

	if err := cs.SubtractBalance(tx.From, tx.Value); err != nil {
		return err
	}

	cs.AddBalance(tx.To, tx.Value)
	return nil
}

func (cs *ChainState) ApplyTransactions(txs []Tx) error {
	for _, tx := range txs {
		if err := cs.ApplyTx(tx); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ChainState) Balances() map[types.Address]types.Amount {
	return cs.balances
}
