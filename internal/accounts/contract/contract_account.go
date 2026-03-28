package contract

import "github.com/chronostech-git/fabrik/internal/types"

type ContractAccount struct {
	balance types.Amount
	alive   bool
	address types.Address
	code    []byte
	storage map[types.Hash]types.Hash
}

func NewAccount(addr types.Address, code []byte) *ContractAccount {
	return &ContractAccount{
		balance: types.ZeroAmount(),
		alive:   true,
		address: addr,
		code:    code,
		storage: make(map[types.Hash]types.Hash),
	}
}

func (ca *ContractAccount) Balance() types.Amount {
	return ca.balance
}

func (ca *ContractAccount) Alive() bool {
	return ca.alive
}

func (ca *ContractAccount) Address() types.Address {
	return ca.address
}

func (ca *ContractAccount) Storage() map[types.Hash]types.Hash {
	return ca.storage
}

func (ca *ContractAccount) UpdateBalance(amount types.Amount) types.Amount {
	newBalance := ca.balance.Add(amount)
	return newBalance
}

func (ca *ContractAccount) Code() []byte {
	return ca.code
}
