package accounts

import (
	"github.com/chronostech-git/fabrik/internal/types"
)

type Account interface {
	Balance() types.Amount
	Alive() bool
	Address() types.Address
	UpdateBalance(amount types.Amount) types.Amount
	Code() []byte
	SetCode(code []byte)
}

func IsContractAccount(account Account) bool {
	return len(account.Code()) > 0
}
