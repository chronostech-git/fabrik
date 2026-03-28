package state

import (
	"errors"

	"github.com/chronostech-git/fabrik/internal/accounts"
	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountState struct {
	Accounts     map[types.Address]accounts.Account
	DeadAccounts int
	Storage      storage.Database
}

func NewAccountState(db storage.Database) *AccountState {
	return &AccountState{
		Accounts:     make(map[types.Address]accounts.Account),
		DeadAccounts: 0,
		Storage:      db,
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

	numDead := len(dead)
	as.DeadAccounts = numDead

	return dead
}

func (as *AccountState) GetAccount(addr types.Address) accounts.Account {
	return as.Accounts[addr]
}

func (as *AccountState) AddAccount(account accounts.Account) {
	as.Accounts[account.Address()] = account
}
