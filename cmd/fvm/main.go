package main

import (
	"log"

	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/chronostech-git/fabrik/internal/storage/memory"
)

func main() {
	state := state.NewAccountState(memory.New())

	instructions, err := fvm.ParseFile("contracts/token.fab")
	if err != nil {
		log.Panic(err)
	}

	bytecode, err := fvm.Compile(instructions)
	if err != nil {
		log.Panic(err)
	}

	deploy := state.Tx{
		From: 1,
		To:   0,
		Data: bytecode,
		Gas:  100000,
	}

	if err := fvm.ApplyTx(state, deploy); err != nil {
		log.Panic(err)
	}

	call := fvm.Tx{
		From: 1,
		To:   1,
		Gas:  100000,
	}

	if err := fvm.ApplyTx(state, call); err != nil {
		log.Panic(err)
	}
}
