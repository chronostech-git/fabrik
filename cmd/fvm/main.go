package main

import (
	"fmt"
	"log"

	"github.com/chronostech-git/fabrik/internal/accounts/external"
	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/keystore"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
	"github.com/chronostech-git/fabrik/internal/types"
)

func main() {
	s := state.NewAccountState(memory.New())

	ks := keystore.NewFileStore("data")

	key, err := ks.GetKey()
	if err != nil {
		log.Panic(err)
	}

	externalAccountAddr := crypto.GenerateAddress(key.PublicKey())
	alice := external.NewAccount(externalAccountAddr)

	s.AddAccount(alice)

	instructions, err := fvm.ParseFile("contracts/complex.fab")
	if err != nil {
		log.Panic(err)
	}

	bytecode, err := fvm.Compile(instructions)
	if err != nil {
		log.Panic(err)
	}

	deploy := &state.Tx{
		From: alice.Address(),
		To:   types.ZeroAddress(),
		Data: bytecode,
		Gas:  100000,
	}

	contractAddr, err := fvm.ApplyTx(s, deploy)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Contract deployed at:", contractAddr)
	fmt.Println(fvm.Disassemble(bytecode))

	call := &state.Tx{
		From: alice.Address(),
		To:   contractAddr,
		Gas:  100000,
	}

	if _, err := fvm.ApplyTx(s, call); err != nil {
		log.Panic(err)
	}

	fmt.Println("Call executed successfully")
}
