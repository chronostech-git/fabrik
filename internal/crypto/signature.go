/*
This file *could* potentially be in the types/ folder as it is technically a type.
For now, it will stay in the crypto/ folder and be used as crypto.Signature or *crypto.Signature (ptr).
*/

package crypto

import (
	"encoding/hex"
	"errors"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s Signature) Bytes() []byte {
	rBytes := s.R.Bytes()
	sBytes := s.S.Bytes()

	out := make([]byte, 64)

	copy(out[32-len(rBytes):32], rBytes)
	copy(out[64-len(sBytes):], sBytes)

	return out
}

func (s Signature) Hex() string {
	return "0x" + hex.EncodeToString(s.Bytes())
}

func BytesToSignature(b []byte) (Signature, error) {
	var sig Signature

	if len(b) != 64 {
		return sig, errors.New("invalid signature length")
	}

	r := new(big.Int).SetBytes(b[:32])
	s := new(big.Int).SetBytes(b[32:])

	sig.R = r
	sig.S = s

	return sig, nil
}
