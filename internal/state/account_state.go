package state

import (
	"errors"

	"github.com/chronostech-git/fabrik/internal/accounts"
	"github.com/chronostech-git/fabrik/internal/accounts/contract"
	"github.com/chronostech-git/fabrik/internal/types"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountState struct {
	Accounts     map[types.Address]accounts.Account
	DeadAccounts int
}

// Create account state
func NewAccountState() *AccountState {
	return &AccountState{
		Accounts:     make(map[types.Address]accounts.Account),
		DeadAccounts: 0,
	}
}

// Update balance given a specific address. If the address is not found, it throws an error.
func (as *AccountState) UpdateBalance(addr types.Address, newBalance types.Amount) error {
	account, ok := as.Accounts[addr]
	if !ok {
		return ErrAccountNotFound
	}
	account.UpdateBalance(newBalance)
	return nil
}

// Not used currently, but will eventually be used to determine how active
// account state is.
func (as *AccountState) FilterDead() map[types.Address]bool {
	dead := make(map[types.Address]bool)
	for addr, account := range as.Accounts {
		if !account.Alive() {
			dead[addr] = true
		}
	}
	as.DeadAccounts = len(dead)
	return dead
}

// NOTE: This may be removed for two reasons -- Lack of use, and simplicity.
func (as *AccountState) GetOrCreateAccount(addr types.Address) accounts.Account {
	acct, ok := as.Accounts[addr]
	if !ok {
		acct = contract.NewAccount(addr)
		as.Accounts[addr] = acct
	}
	return acct
}

// Get an account given an account address
func (as *AccountState) GetAccount(addr types.Address) accounts.Account {
	return as.Accounts[addr]
}

// Add an account to the state
func (as *AccountState) AddAccount(account accounts.Account) {
	as.Accounts[account.Address()] = account
}
