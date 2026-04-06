package blockchain

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

type Wallet struct {
	KeyStore keystore.Store
	Key      *crypto.Key
}

// NewWallet generates a crypto.Key and stores the key using
// the provided keystore.
// Returns a Wallet with the KeyStore and Key filled.
func NewWallet(ks keystore.Store) *Wallet {
	key := crypto.NewKey()
	err := ks.StoreKey(key)
	if err != nil {
		panic(err)
	}

	return &Wallet{
		KeyStore: ks,
		Key:      key,
	}
}

// SignTx uses the crypto.Key generated in NewWallet to
// cryptographically sign a transaction.
func (w *Wallet) SignTx(tx *Transaction) (*crypto.Signature, error) {
	enc, _ := rlp.Encode(tx)
	hash := sha256.Sum256(enc)

	return w.Key.Sign(hash)
}
