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

func NewAccountState() *AccountState {
	return &AccountState{
		Accounts:     make(map[types.Address]accounts.Account),
		DeadAccounts: 0,
	}
}

func (as *AccountState) UpdateBalance(addr types.Address, newBalance types.Amount) error {
	account, ok := as.Accounts[addr]
	if !ok {
		return ErrAccountNotFound
	}
	account.UpdateBalance(newBalance)
	return nil
}

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

func (as *AccountState) GetOrCreateAccount(addr types.Address) accounts.Account {
	acct, ok := as.Accounts[addr]
	if !ok {
		acct = contract.NewAccount(addr)
		as.Accounts[addr] = acct
	}
	return acct
}

func (as *AccountState) GetAccount(addr types.Address) accounts.Account {
	return as.Accounts[addr]
}

func (as *AccountState) AddAccount(account accounts.Account) {
	as.Accounts[account.Address()] = account
}
