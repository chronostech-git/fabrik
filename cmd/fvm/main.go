package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/fvm"
	"github.com/chronostech-git/fabrik/internal/state"
)

var args struct {
	File   string `arg:"--file" help:"Smart contract file (.fab)"`
	Caller string `arg:"--caller" help:"Your key address"`
	Run    string `arg:"--run" help:"Smart contract bytecode hex"`
	Debug  bool   `arg:"--debug" help:"Print debug info"`
	Gas    int    `arg:"--gas" help:"Gas limit" default:"100000"`
}

func runSmartContractFromFile(
	file string,
	accountState *state.AccountState,
	gasLimit uint64,
	debug bool,
) error {
	instructions, err := fvm.ParseFile(file)
	if err != nil {
		return err
	}

	bytecode, err := fvm.Compile(instructions)
	if err != nil {
		return err
	}

	prog := fvm.NewProgram(bytecode)

	vm := fvm.New(prog, accountState, gasLimit)

	if err := vm.Run(); err != nil {
		return err
	}

	if debug {
		fmt.Println("Compiled bytecode:", hex.EncodeToString(bytecode))

		caller, err := hex.DecodeString(args.Caller)
		if err != nil {
			log.Panic(err)
		}

		vm.PrintContractAddress(caller)
		vm.PrintStackData()
		vm.PrintDisasm()
		vm.PrintGasRemaining()
	}

	return nil
}

func runSmartContractFromHex(
	hexCode []byte,
	accountState *state.AccountState,
	gasLimit uint64,
	debug bool,
) error {
	prog := fvm.NewProgram(hexCode)

	vm := fvm.New(prog, accountState, gasLimit)

	if err := vm.Run(); err != nil {
		return err
	}

	if debug {
		fmt.Println("Input (hexidecimal):", hex.EncodeToString(hexCode))
		vm.PrintStackData()
		vm.PrintDisasm()
		vm.PrintGasRemaining()
	}

	return nil
}

func main() {
	arg.MustParse(&args)

	state := state.NewAccountState()

	if args.Run == "" {
		err := runSmartContractFromFile(args.File, state, uint64(args.Gas), args.Debug)
		if err != nil {
			log.Panic(err)
		}
	} else {
		bytecode, err := fvm.HexToBytes(args.Run)
		if err != nil {
			log.Panic(err)
		}
		err = runSmartContractFromHex(bytecode, state, uint64(args.Gas), args.Debug)
		if err != nil {
			log.Panic(err)
		}
	}

	fmt.Println("Smart contract executed.")
}
