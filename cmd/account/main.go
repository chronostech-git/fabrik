package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/accounts/external"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

var args struct {
	DataDir    string
	WithWallet bool
}

func initFileKeyStore(datadir string) *keystore.FileStore {
	return keystore.NewFileStore(datadir)
}

func createNewWallet(datadir string) *blockchain.Wallet {
	ks := initFileKeyStore(datadir)
	return blockchain.NewWallet(ks)
}

func main() {
	arg.MustParse(&args)

	if args.WithWallet {
		wallet := createNewWallet(args.DataDir)
		account := external.NewAccount(wallet.Key.Address)

		log.Printf("An external account with a balance of %s was created\n", account.Balance().String())
		log.Println("A new wallet has been created and was used to create your account")
		log.Printf("Wallet address: %s\n", wallet.Key.Address.String())

		return
	}

	ks := initFileKeyStore(args.DataDir)

	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	account := external.NewAccount(key.Address)

	log.Printf("An external account with a balance of %s was created\n", account.Balance().String())
	log.Printf("Address used for account creation: %s\n", key.Address.String())
}
