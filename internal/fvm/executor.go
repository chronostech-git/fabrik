package fvm

import "github.com/chronostech-git/fabrik/internal/state"

func ApplyTx(state *state.AccountState, tx *state.Tx) error {
	fromAccount := state.GetAccount(tx.From)

	if tx.To.IsZero() {
		// this means it is a contract, so we set the account.Code field
		code := fromAccount.Code()
		tx.Data = code
	}

	to := state.GetAccount(tx.To)

	prog := NewProgram(to.Code())
	vm := New(prog, state, 1000)

	return vm.Run()
}
