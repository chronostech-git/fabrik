package keystore

import "github.com/chronostech-git/fabrik/internal/types"

func FileName(addr types.Address) string {
	return addr.Hex() + ".key"
}
