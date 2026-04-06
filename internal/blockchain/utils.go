package blockchain

import "fmt"

// BlockFilename takes a block and returns the filename used for the disk (.dat file).
// Filename structure: <hash>-<height>.dat
func BlockFilename(b *Block) string {
	return fmt.Sprintf("%s-%08d.dat", b.Hash.String(), b.Header.Height)
}
