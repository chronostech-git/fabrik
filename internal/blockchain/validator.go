package blockchain

type Validator struct {
	chain *Chain // Every validator holds a copy of the chain
	stake *Stake // Minimum stake is 32 fab
}

// Create a new validator
func NewValidator(chain *Chain, stake *Stake) *Validator {
	return &Validator{
		chain: chain,
		stake: stake,
	}
}

// Validate a block using validator
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
