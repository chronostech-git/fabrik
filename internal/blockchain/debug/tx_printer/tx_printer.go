package txprinter

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/chronostech-git/fabrik/internal/blockchain"
)

type TxPrinter struct {
	tx *blockchain.Transaction
}

func New() *TxPrinter {
	return &TxPrinter{
		tx: nil,
	}
}

func (txp *TxPrinter) SetTx(tx *blockchain.Transaction) {
	txp.tx = tx
}

func (txp *TxPrinter) PrintData() {
	if txp.tx == nil {
		log.Panic("Transaction required when calling txprinter.PrintData(...)")
	}

	tx := txp.tx

	fmt.Println("TRANSACTION PRINTER")
	fmt.Printf("Transaction %s\n", tx.Hash.String())
	fmt.Printf("\tSender: %s\n", tx.Sender.String())
	fmt.Printf("\tReceiver: %s\n", tx.Receiver.String())
	fmt.Printf("\tAmount: %s FAB\n", tx.Value.String())

	if tx.Data == nil {
		fmt.Println("\tBytecode: nil")
	}
	fmt.Printf("\tBytecode: %s\n", hex.EncodeToString(tx.Data))

	if tx.Signature == nil {
		fmt.Println("\tSignature: nil")
	}
	fmt.Printf("\tSignature: %s\n", tx.Signature.Hex())
}
