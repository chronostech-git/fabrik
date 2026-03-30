package fvm

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

func ExecuteFromHex(code []byte, accountState *state.AccountState, gasLimit uint64, debug bool) error {
	program := NewProgram(code)
	vm := New(program, accountState, gasLimit)

	err := vm.Run()
	if err != nil {
		return err
	}

	if debug {
		//disasm, err := Disassemble(program.code)
		//if err != nil {
		//	return err
		//}

		vm.PrintStackData()
		vm.PrintGasRemaining()
	}

	return nil
}

func ApplyTx(accountState *state.AccountState, tx *state.Tx) (types.Address, error) {
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

	vm.PrintGasRemaining()
	vm.PrintStackData()

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
