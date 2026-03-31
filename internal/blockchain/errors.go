package blockchain

import "errors"

// Validation
var (
	ErrInvalidBlock        = errors.New("invalid block")
	ErrInvalidTransaction  = errors.New("invalid transaction")
	ErrNoTransactions      = errors.New("no transactions in block")
	ErrInvalidPrevHash     = errors.New("invalid previous hash")
	ErrMaxBlockTimeReached = errors.New("stale block")
)

// Blockchain
var (
	ErrCacheEmpty = errors.New("cache is empty")
)

// Staking
var (
	ErrStakeMinimumNotMet         = errors.New("minimum staking deposit not met")
	ErrInsufficientAccountBalance = errors.New("insufficient account balance for stake")
)
