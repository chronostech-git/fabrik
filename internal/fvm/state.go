package fvm

import (
	"github.com/chronostech-git/fabrik/internal/types"
	"github.com/holiman/uint256"
)

type StateDB interface {
	GetState(addr types.Address, key uint256.Int) uint256.Int
	SetState(addr types.Address, key uint256.Int, value uint256.Int)
}
