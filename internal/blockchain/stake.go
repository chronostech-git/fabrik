package blockchain

import (
	"crypto/sha256"
	"log"

	"github.com/chronostech-git/fabrik/internal/types"
)

var (
	StakeMinimumDeposit = types.NewAmount(32) // Means nothing yet...it does however match Ethereums minimum
)

// A given stake contains a hash, value, and deposit contract address for FVM execution.
type Stake struct {
	hash            types.Hash
	value           types.Amount
	contractAddress types.Address
}

// Initialize a new stake with a given value (must be >= 32 FAB coins)
// NOTE: "32" is arbitrary...this is more than likely gonna change.
func NewStake(value types.Amount) *Stake {
	stake := &Stake{
		value: value,
	}

	stake.hash = sha256.Sum256(stake.value.Bytes())

	return stake
}

// This function is only to be used when a node is marked as
// "dirty".
// A node that has been marked as "dirty", will have been proven to go against
// the consensus in some way shape or form with the intent to write invalid or malicious data.
// This should be done in percentages and strikes, with a 5 strike maximum.
// If the given 5 strike maximum is reached, the validator associated with the stake will have no stake left,
// and will never be able to participate in the network as a validator again.
func (s *Stake) BurnValue(percentage uint64) types.Amount {
	decimal := percentage / 100
	if decimal > 1.0 {
		// greater than 100% of stake is not allowed
		return types.ZeroAmount()
	}

	amountDecimal := types.NewAmount(int64(decimal))

	burnValue := s.value.Mul(amountDecimal) // value * decimal = [amount to burn]

	newStakeValue, err := s.value.Sub(burnValue) // subtract [amount to burn] from the current stake value and return the new value
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
