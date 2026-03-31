package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/accounts"
	"github.com/chronostech-git/fabrik/internal/accounts/contract"
	"github.com/chronostech-git/fabrik/internal/accounts/external"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

var args struct {
	DataDir    string `arg:"required"`
	Type       string `arg:"--type,required"`
	WithWallet bool   `arg:"--wallet"`
}

func main() {
	arg.MustParse(&args)

	state := state.NewAccountState()

	var account accounts.Account
	var store keystore.Store
	var key *crypto.Key
	var wallet *blockchain.Wallet

	store = keystore.NewFileStore(args.DataDir)

	if args.WithWallet {
		wallet = blockchain.NewWallet(store)
		key = wallet.Key
		log.Printf("Wallet created address=%s", key.Address.String())
	} else {
		var err error
		key, err = store.GetKey()
		if err != nil {
			log.Panicf("failed to load key from datadir: %v", err)
		}
		log.Printf("Wallet loaded from datadir=%s address=%s", args.DataDir, key.Address.String())
	}

	switch args.Type {
	case "contract":
		account = contract.NewAccount(key.Address)
	case "external":
		account = external.NewAccount(key.Address)
	default:
		log.Panicf("unknown account type: %s", args.Type)
	}

	state.AddAccount(account)

	log.Printf("%s account created using wallet %s", strings.ToUpper(args.Type), key.Address.String())

	fmt.Println()
	if wallet != nil {
		log.Println("Public key:", wallet.Key.PublicKeyHex())
	}
	log.Printf("Private key stored on disk: %s/keystore/<address>.key", args.DataDir)
	fmt.Println()
	log.Println("WARN: Do not share your secure private key with anyone.")
}
