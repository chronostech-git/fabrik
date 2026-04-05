package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir string

	// Transaction args
	Transact bool   `arg:"--transact" help:"Use this flag if you would like to initiate a transaction"`
	SendTo   string `arg:"--sendto" help:"Receiver (address) of transaction"`
	Amount   int    `arg:"--amount" help:"Amount in FAB coin"`
	GasLimit int    `arg:"--gaslimit" help:"Gas limit for transaction"`

	// Private key signing args
	Sign      bool   `arg:"--sign" help:"Use this flag to sign data using private key."`
	Data      string `arg:"--data" help:"Message to sign (string)"`
	Verify    bool   `arg:"--verify" help:"Use this flag to verify the signature generated when signing data"`
	Signature string `arg:"--sig" help:"Signature you would like to verify"`
}

func createTransaction(datadir string, sentTo string, amount int, gasLimit int) *blockchain.Transaction {
	toAddr, err := types.HexToAddress(sentTo)
	if err != nil {
		log.Panic(err)
	}

	amountToSend := types.NewAmount(int64(amount))

	keystore := keystore.NewFileStore(datadir)

	key, err := keystore.GetKey()
	if err != nil {
		log.Panic(err)
	}

	return blockchain.NewTx(key.Address, toAddr, amountToSend, 0, nil)
}

func signData(datadir string, data string, verify bool) *crypto.Signature {
	keystore := keystore.NewFileStore(datadir)

	key, err := keystore.GetKey()
	if err != nil {
		log.Panic("failed to retreive key from datadir")
	}

	dataHash := sha256.Sum256([]byte(data))

	sig, err := key.Sign(dataHash)
	if err != nil {
		log.Panic("failed to sign hashed data")
	}

	if verify {
		validSig := key.Verify(dataHash, sig)
		if !validSig {
			log.Panic("failed to verify signature: invalid signature")
		}
		return sig
	}

	return sig
}

func main() {
	arg.MustParse(&args)

	if args.Transact {
		tx := createTransaction(args.DataDir, args.SendTo, args.Amount, args.GasLimit)
		log.Printf("transaction of %s FAB sent from %s to %s\n", tx.Value.String(), tx.Sender.String(), tx.Receiver.String())
	} else if args.Sign {
		sig := signData(args.DataDir, args.Data, args.Verify)
		fmt.Printf("signed & verified signature: %s", sig.Hex())
	}
}
