package blockchain

import "errors"

// Validation
var (
	ErrInvalidBlock       = errors.New("invalid block")
	ErrInvalidTransaction = errors.New("invalid transaction")
)

// Blockchain
var (
	ErrCacheEmpty = errors.New("cache is empty")
)
