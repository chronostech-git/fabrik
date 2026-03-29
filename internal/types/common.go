package types

import (
	"encoding/hex"
	"errors"
)

type Hash [32]byte

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) String() string {
	return h.Hex()
}

func (h Hash) IsZero() bool {
	var zero Hash
	return h == zero
}

func BytesToHash(b []byte) Hash {
	var h Hash
	if len(b) > 32 {
		b = b[len(b)-32:]
	}
	copy(h[32-len(b):], b)
	return h
}

func HexToHash(s string) (Hash, error) {
	var h Hash

	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return h, err
	}

	if len(b) != 32 {
		return h, errors.New("invalid hash length")
	}

	copy(h[:], b)
	return h, nil
}

func Empty32() Hash {
	return Hash{}
}

type Address [20]byte

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) String() string {
	return a.Hex()
}

func (a Address) IsZero() bool {
	var zero Address
	return a == zero
}

func BytesToAddress(b []byte) Address {
	var a Address
	if len(b) > 20 {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)
	return a
}

func HexToAddress(s string) (Address, error) {
	var a Address

	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return a, err
	}

	if len(b) != 20 {
		return a, errors.New("invalid address length")
	}

	copy(a[:], b)
	return a, nil
}

func ZeroAddress() Address {
	var addr Address
	return addr
}
