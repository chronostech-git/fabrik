package state

import "github.com/chronostech-git/fabrik/internal/types"

type BlockView interface {
	Transactions() []TxView
}

type TxView interface {
	From() types.Address
	To() types.Address
	Val() types.Amount
}
