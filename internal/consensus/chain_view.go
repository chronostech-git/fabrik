package consensus

import (
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage"
)

// ChainView is a consensus compatible version of the
// blockchain. We take what we need, and that's it.
type ChainView struct {
	DB    storage.Database
	State *state.ChainState
	Head  *BlockView
}
