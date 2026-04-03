package blockchain

import (
	"encoding/json"
	"log"

	"github.com/chronostech-git/fabrik/internal/accounts/contract"
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/types"
)

const (
	DepositContractAddress = "0x8114dc1018aaacdbb788f9a6f58d4460234188cb"
)

type StakeDepositReceipt struct {
	Value           int64
	Hash            string
	GasUsed         int64
	ContractAddress string
	TxSig           string
}

func NewStakeDepositReceipt(depositTx *Transaction, gasUsed int64, contractAddress string) *StakeDepositReceipt {
	return &StakeDepositReceipt{
		Value:           depositTx.Value.Big().Int64(),
		Hash:            depositTx.Hash.String(),
		GasUsed:         gasUsed,
		ContractAddress: contractAddress,
		TxSig:           depositTx.Signature.Hex(),
	}
}

func (dr *StakeDepositReceipt) Json() string {
	j, err := json.MarshalIndent(dr, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	return string(j)
}

func CreateStakeDepositTransaction(
	sender types.Address,
	amount types.Amount,
	accountState *state.AccountState,
	gasLimit uint64,
	code []byte,
	debug bool,
) (*Transaction, uint64, error) {
	hexToAddr, err := types.HexToAddress(DepositContractAddress)
	if err != nil {
		return nil, 0, err
	}

	stake := NewStake(amount)
	stake.contractAddress = hexToAddr

	depositTx := NewTx(sender, hexToAddr, stake.value, 0, code)

	fvmProg := fvm.NewProgram(code)

	vm := fvm.New(fvmProg, accountState, gasLimit)
	if err := vm.Run(); err != nil {
		return nil, 0, err
	}

	if accountState.GetAccount(hexToAddr) == nil {
		depositAccount := contract.NewAccount(hexToAddr)
		accountState.AddAccount(depositAccount)
	}

	err = accountState.UpdateBalance(hexToAddr, amount)
	if err != nil {
		return nil, 0, err
	}

	if debug {
		vm.PrintStackData()
		vm.PrintDisasm()
	}

	return depositTx, vm.GasRemaining(), nil
}
