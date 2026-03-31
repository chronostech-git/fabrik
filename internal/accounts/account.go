package accounts

import (
	"github.com/chronostech-git/fabrik/internal/types"
)

type Account interface {
	Balance() types.Amount
	Alive() bool
	Address() types.Address
	UpdateBalance(amount types.Amount)
	Code() []byte        // For contract accounts only
	SetCode(code []byte) // For contract accounts only
}

// Determine if a given account is a contract account by
// the length of account.Code. If code is present, it is a contract account.
func IsContractAccount(account Account) bool {
	return len(account.Code()) > 0
}
