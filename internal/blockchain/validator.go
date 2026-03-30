package blockchain

import "github.com/chronostech-git/fabrik/internal/state"

type Validator struct {
	chain *Chain
	state *state.ChainState
}
