package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

func BytesToPrivateKey(b []byte) *ecdsa.PrivateKey {
	curve := elliptic.P256()

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve

	priv.D = new(big.Int).SetBytes(b)

	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(b)

	return priv
}
