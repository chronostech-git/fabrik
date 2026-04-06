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
	Data  []byte // For contract execution.
	Gas   uint64
}

type ChainState struct {
	balances map[types.Address]types.Amount
}

// NewChainState creates a new instance of ChainState
// with an empty map of balances
func NewChainState() *ChainState {
	return &ChainState{
		balances: make(map[types.Address]types.Amount),
	}
}

// GetBalance get a balance assuming a correct and valid address has been provided.
func (cs *ChainState) GetBalance(addr types.Address) types.Amount {
	bal, ok := cs.balances[addr]
	if !ok {
		return types.NewAmount(0)
	}
	return bal
}

// SetBalance set a balance using an account address and amount
func (cs *ChainState) SetBalance(addr types.Address, amount types.Amount) {
	cs.balances[addr] = amount
}

// AddBalance add to current balance of a given address
func (cs *ChainState) AddBalance(addr types.Address, amount types.Amount) {
	current := cs.GetBalance(addr)
	cs.balances[addr] = current.Add(amount)
}

// SubtractBalance subtract from current balance
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

// ApplyTx change state based on specified transaction
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

// ApplyTransactions changes the current state based on multiple transactions provided.
func (cs *ChainState) ApplyTransactions(txs []Tx) error {
	for _, tx := range txs {
		if err := cs.ApplyTx(tx); err != nil {
			return err
		}
	}
	return nil
}

// Balances gets all balances (returns map[<account-address>]<amount>)
func (cs *ChainState) Balances() map[types.Address]types.Amount {
	return cs.balances
}
