package state

import "github.com/chronostech-git/fabrik/internal/types"

type TxView interface {
	From() types.Address
	To() types.Address
	Val() types.Amount // Transaction value
	Dat() []byte       // Contract bytecode
}
