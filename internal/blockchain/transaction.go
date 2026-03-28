package blockchain

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/types"
)

type Transaction struct {
	Sender   types.Address
	Receiver types.Address
	Value    types.Amount
	Nonce    uint64
	Data     []byte

	Signature *crypto.Signature
	Hash      types.Hash
}

func NewTx(sender, receiver types.Address, value types.Amount, nonce uint64, data []byte) *Transaction {
	tx := &Transaction{
		Sender:   sender,
		Receiver: receiver,
		Value:    value,
		Nonce:    nonce,
		Data:     data,
	}

	tx.Hash = tx.computeHash()
	return tx
}

func (tx *Transaction) computeHash() types.Hash {
	type txData struct {
		Sender   types.Address
		Receiver types.Address
		Value    types.Amount
		Nonce    uint64
		Data     []byte
	}

	d := txData{
		Sender:   tx.Sender,
		Receiver: tx.Receiver,
		Value:    tx.Value,
		Nonce:    tx.Nonce,
		Data:     tx.Data,
	}

	enc, _ := rlp.Encode(d)
	return sha256.Sum256(enc)
}

// TxView interface functions
func (tx *Transaction) From() types.Address {
	return tx.Sender
}

func (tx *Transaction) To() types.Address {
	return tx.Receiver
}

func (tx *Transaction) Val() types.Amount {
	return tx.Value
}

func (tx *Transaction) Dat() []byte {
	return tx.Data
}
