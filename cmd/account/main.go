package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/accounts"
	"github.com/chronostech-git/fabrik/internal/accounts/contract"
	"github.com/chronostech-git/fabrik/internal/accounts/external"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
)

var args struct {
	DataDir string `arg:"required"`
	Type    string `arg:"--type,required"`
}

func main() {
	arg.MustParse(&args)

	state := state.NewAccountState()
	keystore := keystore.NewFileStore(args.DataDir)

	var account accounts.Account

	key, err := keystore.GetKey()
	if err != nil {
		log.Panic(err)
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
	fmt.Printf("%s account created using wallet %s", strings.ToUpper(args.Type), key.Address.String())
}
