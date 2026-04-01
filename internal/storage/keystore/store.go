package keystore

import "github.com/chronostech-git/fabrik/internal/crypto"

// Basic for file storage operations / functionality
// This can however be extended to support multiple ways of
// storing a secret key (private key).
type Store interface {
	GetKey() (*crypto.Key, error)
	StoreKey(*crypto.Key) error
}
