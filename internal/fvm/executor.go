package fvm

import (
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

// When the FVM runs, things change in the state.
// We can apply the transaction to the state so long as the Virtual Machine
// runs beginning to end.
func ApplyTx(accountState *state.AccountState, tx *state.Tx, debug bool) (types.Address, error) {
	if tx.To.IsZero() {
		contractAddr, err := DeriveContractAddress(tx.From.Bytes(), 0)
		if err != nil {
			return types.ZeroAddress(), err
		}
		addr := types.BytesToAddress(contractAddr)

		contractAccount := accountState.GetOrCreateAccount(addr)
		contractAccount.SetCode(tx.Data)
		accountState.AddAccount(contractAccount)

		prog := NewProgram(tx.Data)
		vm := New(prog, accountState, tx.Gas)
		if err := vm.Run(); err != nil {
			return types.ZeroAddress(), err
		}
		return addr, nil
	}

	contractAccount := accountState.GetAccount(tx.To)
	code := contractAccount.Code()

	prog := NewProgram(code)
	vm := New(prog, accountState, tx.Gas)
	if err := vm.Run(); err != nil {
		return types.Address{}, err
	}

	if debug {
		vm.PrintGasRemaining()
		vm.PrintStackData()
	}

	return tx.To, nil
}
