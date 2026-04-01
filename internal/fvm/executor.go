package fvm

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

// When the FVM runs, things change in the state.
// We can apply the transaction to the state so long as the Virtual Machine
// runs beginning to end.
func ApplyTx(accountState *state.AccountState, tx *state.Tx, debug bool) (types.Address, error) {
	if tx.To.IsZero() {
		contractAddr := deriveContractAddress(tx.From, tx.Data)

		contractAccount := accountState.GetOrCreateAccount(contractAddr)
		contractAccount.SetCode(tx.Data)
		accountState.AddAccount(contractAccount)

		prog := NewProgram(tx.Data)
		vm := New(prog, accountState, tx.Gas)
		if err := vm.Run(); err != nil {
			return types.Address{}, err
		}

		return contractAddr, nil
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

func deriveContractAddress(deployer types.Address, code []byte) types.Address {
	h := sha256.New()
	h.Write(deployer[:])
	h.Write(code)
	sum := h.Sum(nil)

	var addr types.Address
	copy(addr[:], sum[len(sum)-len(addr):])
	return addr
}
