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
	Type       string `arg:"--type,required"` // external or contract
	WithWallet bool   `arg:"--wallet"`        // If true, a new wallet will be created. Otherwise, loadWalletAndKey is used.
}

func createNewWallet(ks keystore.Store) (*blockchain.Wallet, *crypto.Key) {
	wallet := blockchain.NewWallet(ks)
	key := wallet.Key
	return wallet, key
}

func loadWalletAndKey(ks keystore.Store) (*blockchain.Wallet, *crypto.Key) {
	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}
	return &blockchain.Wallet{
		KeyStore: ks,
		Key:      key,
	}, nil
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
		wallet, key = createNewWallet(store)
	} else {
		wallet, key = loadWalletAndKey(store)
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
