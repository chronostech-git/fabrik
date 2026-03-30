package blockchain

import "fmt"

func BlockFilename(b *Block) string {
	return fmt.Sprintf("%s-%08d.dat", b.Hash.String(), b.Header.Height)
}
