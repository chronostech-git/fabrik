package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

// Takes in the bytes of an encoded private key and returns type *ecdsa.PrivateKey
func BytesToPrivateKey(b []byte) *ecdsa.PrivateKey {
	d := new(big.Int).SetBytes(b)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = d
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d.Bytes())
	return priv
}
