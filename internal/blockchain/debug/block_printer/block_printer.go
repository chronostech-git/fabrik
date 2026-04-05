package blockprinter

import (
	"fmt"
	"log"

	"github.com/chronostech-git/fabrik/internal/blockchain"
)

type BlockPrinter struct {
	block *blockchain.Block
}

func New() *BlockPrinter {
	return &BlockPrinter{
		block: nil,
	}
}

func (bp *BlockPrinter) SetBlock(b *blockchain.Block) {
	bp.block = b
}

func (bp *BlockPrinter) PrintData() {
	if bp.block == nil {
		log.Panic("Block required when calling blockprinter.PrintData(...)")
	}

	b := bp.block

	fmt.Println("BLOCK PRINTER")
	fmt.Printf("Block %s\n", b.Hash.String())
	fmt.Printf("\tPrevious hash: %s\n", b.Header.PrevHash.String())
	fmt.Printf("\tTimestamp: %d\n", b.Header.Timestamp)
	fmt.Printf("\tRoot: %s\n", b.Header.TxRoot.String())
	fmt.Printf("\tHeight: %08d\n", b.Header.Height)
	fmt.Println("Tx hashes")
	// TODO Print GasLimit, GasUsed, and the BaseFee (per gas)

	totalHashes := 0
	for _, tx := range b.Txs {
		fmt.Printf("\tHash #%d: %s\n", totalHashes+1, tx.Hash.String())
		totalHashes += 1
	}
}
