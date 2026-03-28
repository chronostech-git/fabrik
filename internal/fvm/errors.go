package fvm

import "errors"

var (
	ErrInvalidOpcode = errors.New("invalid opcode")
	ErrOutOfBounds   = errors.New("out of bounds")
)
