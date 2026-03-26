package blockchain

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/accounts/external"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

type Wallet struct {
	KeyStore keystore.Store
	Key      *crypto.Key
}

func NewWallet(ks keystore.Store) *Wallet {
	key := crypto.NewKey()
	_ = ks.StoreKey(key)

	return &Wallet{
		KeyStore: ks,
		Key:      key,
	}
}

func (w *Wallet) CreateExternalAccount() *external.ExternalAccount {
	return external.NewAccount(w.Key.Address)
}

func (w *Wallet) SignTx(tx *Transaction) (*crypto.Signature, error) {
	enc, _ := rlp.Encode(tx)
	hash := sha256.Sum256(enc)

	return w.Key.Sign(hash)
}
