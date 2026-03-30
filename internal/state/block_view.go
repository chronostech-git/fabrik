package state

type BlockView interface {
	Transactions() []TxView
}
