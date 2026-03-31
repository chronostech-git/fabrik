package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
)

var args struct {
	File  string `arg:"--file" help:"Smart contract file (.fab)"`
	Debug bool   `arg:"--debug" help:"Print debug info"`
	Gas   int    `arg:"--gas" help:"Gas limit" default:"100000"`
}

func main() {
	arg.MustParse(&args)

	state := state.NewAccountState()

	instructions, err := fvm.ParseFile(args.File)
	if err != nil {
		log.Panic(err)
	}

	bytecode, err := fvm.Compile(instructions)
	if err != nil {
		log.Panic(err)
	}

	program := fvm.NewProgram(bytecode)

	vm := fvm.New(program, state, uint64(args.Gas))

	if err := vm.Run(); err != nil {
		log.Panic(err)
	}

	fmt.Println()
	if args.Debug {
		vm.PrintStackData()
		vm.PrintGasRemaining()
	}
}
