package external

import (
	"github.com/chronostech-git/fabrik/internal/types"
)

type ExternalAccount struct {
	balance types.Amount
	alive   bool
	address types.Address
}

func NewAccount(addr types.Address) *ExternalAccount {
	return &ExternalAccount{
		balance: types.ZeroAmount(),
		alive:   true,
		address: addr,
	}
}

func (ea *ExternalAccount) Balance() types.Amount {
	return ea.balance
}

func (ea *ExternalAccount) Alive() bool {
	return ea.alive
}

func (ea *ExternalAccount) Address() types.Address {
	return ea.address
}

func (ea *ExternalAccount) UpdateBalance(amount types.Amount) types.Amount {
	newBalance := ea.balance.Add(amount)
	return newBalance
}

func (ea *ExternalAccount) Code() []byte {
	return nil
}
