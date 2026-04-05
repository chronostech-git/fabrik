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
	Type       string `arg:"--type,required"`     // external or contract
	WithWallet bool   `arg:"--wallet"`            // Create a new wallet if true
	Stake      int    `arg:"--stake" default:"0"` // Stake amount (>=32 fab)
	GasLimit   int    `arg:"--gas" help:"Gas limit for calling staking deposit contract"`
	Debug      bool   `arg:"--debug" help:"Debug mode for FVM deposit contract execution"`
}

func createWallet(store keystore.Store) (*blockchain.Wallet, *crypto.Key) {
	w := blockchain.NewWallet(store)
	return w, w.Key
}

func loadWallet(store keystore.Store) (*blockchain.Wallet, *crypto.Key) {
	key, err := store.GetKey()
	if err != nil {
		log.Panic(err)
	}
	return &blockchain.Wallet{KeyStore: store, Key: key}, key
}

func main() {
	arg.MustParse(&args)

	store := keystore.NewFileStore(args.DataDir)
	wallet, key := func() (*blockchain.Wallet, *crypto.Key) {
		if args.WithWallet {
			return createWallet(store)
		}
		return loadWallet(store)
	}()

	state := state.NewAccountState()

	var account accounts.Account
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
	if wallet != nil {
		log.Println("Public key:", wallet.Key.PublicKeyHex())
	}
	log.Printf("Private key stored on disk: %s/keystore/<address>.key", args.DataDir)
	fmt.Println()
	log.Println("WARN: Do not share your secure private key with anyone.")
}
