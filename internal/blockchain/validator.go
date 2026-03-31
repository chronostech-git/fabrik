package blockchain

import (
	"github.com/chronostech-git/fabrik/internal/accounts/contract"
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

type Validator struct {
	chain   *Chain // Every validator
	stake   *Stake
	account *contract.ContractAccount
}

func NewValidator(chain *Chain, account *contract.ContractAccount) *Validator {
	return &Validator{
		chain:   chain,
		account: account,
	}
}

func (v *Validator) SetStake(amount types.Amount) error {
	if amount.Cmp(StakeMinimumDeposit) < 32 {
		return ErrStakeMinimumNotMet
	}

	s := NewStake(amount)
	s.contractAddress = v.account.Address()

	v.account.UpdateBalance(s.value)

	return nil
}

func (v *Validator) executeDepositContractCode(account *contract.ContractAccount) error {
	instructions, err := fvm.ParseFile("/var/home/caleb/dev/fabrik/contracts/deposit.fab")
	if err != nil {
		return err
	}

	bytecode, err := fvm.Compile(instructions)
	if err != nil {
		return err
	}

	account.SetCode(bytecode)

	state := state.NewAccountState()

	prog := fvm.NewProgram(bytecode)
	vm := fvm.New(prog, state, 100)

	if err := vm.Run(); err != nil {
		return err
	}

	vm.PrintGasRemaining()
	vm.PrintStackData()

	return nil
}

func (v *Validator) ValidateBlock(b *Block) (bool, error) {
	if v.chain.Genesis == nil {
		return false, ErrGenesisMissing
	}
	if len(b.Txs) == 0 {
		return false, ErrNoTransactions
	}
	if b.Header.Timestamp > b.Header.Timestamp+int64(MaxFutureBlockTime) {
		return false, ErrMaxBlockTimeReached
	}

	return true, nil
}
