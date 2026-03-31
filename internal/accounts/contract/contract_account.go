package contract

import "github.com/chronostech-git/fabrik/internal/types"

type ContractAccount struct {
	balance types.Amount
	alive   bool
	address types.Address
	code    []byte
	storage map[types.Hash]types.Hash
}

// Create a new contract account
func NewAccount(addr types.Address) *ContractAccount {
	return &ContractAccount{
		balance: types.ZeroAmount(),
		alive:   true,
		address: addr,
		storage: make(map[types.Hash]types.Hash),
	}
}

func (ca *ContractAccount) Code() []byte {
	return ca.code
}

func (ca *ContractAccount) SetCode(code []byte) {
	ca.code = code
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

// Update the contract account balance
func (ca *ContractAccount) UpdateBalance(amount types.Amount) {
	newBalance := ca.balance.Add(amount)
	ca.balance = newBalance
}
