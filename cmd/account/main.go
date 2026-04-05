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
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/types"
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

	// Handle staking deposit if requested
	if args.Stake > 0 {
		codeBytes, err := fvm.HexToBytes("333455424400")
		if err != nil {
			log.Panic(err)
		}

		stakeAmount := types.NewAmount(int64(args.Stake))
		depositTx, gasRemaining, err := blockchain.CreateStakeDepositTransaction(
			account.Address(),
			stakeAmount,
			state,
			uint64(args.GasLimit),
			codeBytes,
			args.Debug,
		)
		if err != nil {
			log.Panic(err)
		}

		sig, err := key.Sign(depositTx.Hash)
		if err != nil {
			log.Panic(err)
		}
		depositTx.Signature = sig

		gasUsed := args.GasLimit - int(gasRemaining)
		stakeReceipt := blockchain.NewStakeDepositReceipt(depositTx, int64(gasUsed), blockchain.DepositContractAddress)

		log.Println("NOTE: This account is now allowed to become a validator. Below is your custom Stake Receipt:")
		fmt.Println(stakeReceipt.Json())
	}

	log.Printf("%s account created using wallet %s", strings.ToUpper(args.Type), key.Address.String())
	if wallet != nil {
		log.Println("Public key:", wallet.Key.PublicKeyHex())
	}
	log.Printf("Private key stored on disk: %s/keystore/<address>.key", args.DataDir)
	fmt.Println()
	log.Println("WARN: Do not share your secure private key with anyone.")
}
