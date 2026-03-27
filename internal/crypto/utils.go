package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

func BytesToPrivateKey(b []byte) *ecdsa.PrivateKey {
	d := new(big.Int).SetBytes(b)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = d
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d.Bytes())
	return priv
}
