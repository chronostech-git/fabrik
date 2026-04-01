package crypto

import (
	"crypto/ecdsa"

	"github.com/chronostech-git/fabrik/internal/types"
	"golang.org/x/crypto/sha3"
)

// Generates a 20-bit address derived from a cryptographic public key
func GenerateAddress(pub *ecdsa.PublicKey) types.Address {
	x := pub.X.Bytes()
	y := pub.Y.Bytes()

	xPadded := make([]byte, 32)
	yPadded := make([]byte, 32)

	copy(xPadded[32-len(x):], x)
	copy(yPadded[32-len(y):], y)

	pubBytes := append(xPadded, yPadded...)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes)
	sum := hash.Sum(nil)

	return types.BytesToAddress(sum[len(sum)-20:])
}
