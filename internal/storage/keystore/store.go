package keystore

import "github.com/chronostech-git/fabrik/internal/crypto"

type Store interface {
	GetKey() (*crypto.Key, error)
	StoreKey(*crypto.Key) error
}
