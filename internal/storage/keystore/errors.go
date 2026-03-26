package keystore

import "errors"

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrInvalidKey      = errors.New("invalid key")
	ErrAddressMismatch = errors.New("address mismatch")
)
