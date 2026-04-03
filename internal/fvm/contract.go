package fvm

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// DeriveContractAddress deterministically computes a contract address
// given the creator's address (20 bytes) and their nonce.
func DeriveContractAddress(creator []byte, nonce uint64) ([]byte, error) {
	if len(creator) != 20 {
		return nil, fmt.Errorf("creator address must be 20 bytes")
	}

	// Concatenate creator address + nonce (big endian 8 bytes)
	data := make([]byte, 0, 28)
	data = append(data, creator...)

	nonceBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		nonceBytes[7-i] = byte(nonce >> (8 * i))
	}
	data = append(data, nonceBytes...)

	// Hash the data
	hash := sha256.Sum256(data)

	contractAddr := hash[12:]

	return contractAddr, nil
}

// Helper to return as hex string
func DeriveContractAddressHex(creator []byte, nonce uint64) (string, error) {
	addr, err := DeriveContractAddress(creator, nonce)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(addr), nil
}
