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
	Type       string `arg:"--type,required"`                                             // external or contract
	WithWallet bool   `arg:"--wallet"`                                                    // If true, a new wallet will be created. Otherwise, loadWalletAndKey is used.
	Stake      int    `arg:"--stake" default:"0"`                                         // If present, it will set the stake amount (must be >= 32 fab)
	Gas        int    `arg:"--gas" help:"Gas limit for calling staking deposit contract"` // Gas limit for calling staking deposit contract
	Debug      bool   `arg:"--debug" help:"Debug mode for FVM deposit contract execution"`
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

	if args.Stake != 0 {
		codeHexToBytes, err := fvm.HexToBytes("333455424400")
		if err != nil {
			log.Panic(err)
		}

		stakeToAmount := types.NewAmount(int64(args.Stake))
		depositTx, gasRemaining, err := blockchain.CreateStakeDepositTransaction(
			account.Address(),
			stakeToAmount,
			state,
			uint64(args.Gas),
			codeHexToBytes,
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

		gasUsed := args.Gas - int(gasRemaining)

		stakeReceipt := blockchain.NewStakeDepositReceipt(
			depositTx,
			int64(gasUsed),
			blockchain.DepositContractAddress,
		)
		log.Printf("NOTE: This account is now allowed to become a validator. Below is your custom Stake Receipt.")
		fmt.Println()
		fmt.Println(stakeReceipt.Json())
		fmt.Println()
	}

	log.Printf("%s account created using wallet %s", strings.ToUpper(args.Type), key.Address.String())

	fmt.Println()
	if wallet != nil {
		log.Println("Public key:", wallet.Key.PublicKeyHex())
	}
	log.Printf("Private key stored on disk: %s/keystore/<address>.key", args.DataDir)
	fmt.Println()
	log.Println("WARN: Do not share your secure private key with anyone.")
}
