package blockchain

import "errors"

// Blockchain related errors
var (
	ErrCacheEmpty  = errors.New("cache is empty")
	ErrFailedWrite = errors.New("write block failed")
)
