package blockchain

import (
	"crypto/sha256"
	"log"

	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	StakeMinimumDeposit = types.NewAmount(32)
)

type Stake struct {
	hash            types.Hash
	value           types.Amount
	contractAddress types.Address
}

func NewStake(value types.Amount) *Stake {
	stake := &Stake{
		value: value,
	}

	stake.hash = sha256.Sum256(stake.value.Bytes())

	return stake
}

func (s *Stake) BurnValue(percentage uint64) types.Amount {
	decimal := percentage / 100
	if decimal > 1.0 {
		// greater than 100% of stake is not allowed
		return types.ZeroAmount()
	}

	amountDecimal := types.NewAmount(int64(decimal))

	burnValue := s.value.Mul(amountDecimal) // value * decimal = [amount to burn]

	newStakeValue, err := s.value.Sub(burnValue)
	if err != nil {
		log.Panic(err)
	}

	return newStakeValue
}

func (s *Stake) Hash() types.Hash {
	return s.hash
}

func (s *Stake) Value() types.Amount {
	return s.value
}

func (s *Stake) ContractAddress() types.Address {
	return s.contractAddress
}
