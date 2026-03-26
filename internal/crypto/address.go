package crypto

import (
	"crypto/ecdsa"

	"github.com/chronostech-git/fabrik/internal/types"
	"golang.org/x/crypto/sha3"
)

func GenerateAddress(pub *ecdsa.PublicKey) types.Address {
	pubBytes := append(pub.X.Bytes(), pub.Y.Bytes()...)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes)
	sum := hash.Sum(nil)

	return types.BytesToAddress(sum[len(sum)-20:])
}
