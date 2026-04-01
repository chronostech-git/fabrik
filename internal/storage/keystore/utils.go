package keystore

import "github.com/chronostech-git/fabrik/internal/types"

// Used to create a .key file in file_store.go
func FileName(addr types.Address) string {
	return addr.Hex() + ".key"
}
