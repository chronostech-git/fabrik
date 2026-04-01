package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"

	"github.com/chronostech-git/fabrik/internal/types"
)

type Key struct {
	PrivateKey *ecdsa.PrivateKey
	Address    types.Address
}

// Create a new key
func NewKey() *Key {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	k := &Key{
		PrivateKey: priv,
	}

	k.Address = GenerateAddress(&priv.PublicKey)
	return k
}

// Sign hashed data using the ecdsa private key
func (k *Key) Sign(hash types.Hash) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return &Signature{R: r, S: s}, nil
}

// Verify if a signature is valid given the data's hash, and signature itself
func (k *Key) Verify(hash types.Hash, sig *Signature) bool {
	pub := &k.PrivateKey.PublicKey
	return ecdsa.Verify(pub, hash[:], sig.R, sig.S)
}

func (k *Key) PublicKey() *ecdsa.PublicKey {
	return &k.PrivateKey.PublicKey
}

func (k *Key) PublicKeyBytes() []byte {
	publicKeyBytes := append(k.PublicKey().X.Bytes(), k.PublicKey().Y.Bytes()...)
	return publicKeyBytes
}

func (k *Key) PublicKeyHex() string {
	return hex.EncodeToString(k.PublicKeyBytes())
}
