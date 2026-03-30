package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/blockchain"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
	"github.com/chronostech-git/fabrik/internal/types"
)

var args struct {
	DataDir      string `arg:"required"`
	TestMode     bool   `arg:"--testwallet"`
	StartBalance int64  `arg:"--balance"`
}

func main() {
	arg.MustParse(&args)

	state := state.NewAccountState(memory.New())

	if err := os.MkdirAll(args.DataDir, 0700); err != nil {
		log.Panic(err)
	}

	ks := keystore.NewFileStore(args.DataDir)

	wallet := blockchain.NewWallet(ks)
	account := state.GetOrCreateAccount(wallet.Key.Address)

	if args.TestMode {
		startBalance := account.UpdateBalance(types.NewAmount(args.StartBalance))
		state.SetAccount(account)
		state.UpdateBalance(account.Address(), types.NewAmount(args.StartBalance))
		log.Printf("Created test wallet with starting balance of %s", startBalance.String())
	}

	fmt.Println("Address:", account.Address().Hex())

	if args.TestMode {
		fmt.Println("Balance:", account.Balance().String())
	}
}
