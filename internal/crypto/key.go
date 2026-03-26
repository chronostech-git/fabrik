package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/chronostech-git/fabrik/internal/types"
)

type Key struct {
	PrivateKey *ecdsa.PrivateKey
	Address    types.Address
}

func NewKey() *Key {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	k := &Key{
		PrivateKey: priv,
	}

	k.Address = GenerateAddress(&priv.PublicKey)
	return k
}

func (k *Key) Sign(hash types.Hash) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return &Signature{R: r, S: s}, nil
}

func (k *Key) Verify(hash types.Hash, sig *Signature) bool {
	pub := &k.PrivateKey.PublicKey
	return ecdsa.Verify(pub, hash[:], sig.R, sig.S)
}

func (k *Key) PublicKey() *ecdsa.PublicKey {
	return &k.PrivateKey.PublicKey
}
