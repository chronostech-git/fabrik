package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chronostech-git/fabrik/internal/fvm"
)

var args struct {
	File   string `arg:"--file" help:"Smart contract file (.fab)"`
	Caller string `arg:"--caller" help:"Your key address (hex)"`
	Run    string `arg:"--run" help:"Smart contract bytecode hex"`
	Debug  bool   `arg:"--debug" help:"Print debug info"`
	Gas    int    `arg:"--gas" help:"Gas limit" default:"100000"`
}

func main() {
	arg.MustParse(&args)

	var bytecode []byte
	var err error

	// Compile from file if --run not provided
	if args.Run == "" {
		instructions, err := fvm.ParseFile(args.File)
		if err != nil {
			log.Panicf("failed to parse file: %v", err)
		}

		bytecode, err = fvm.Compile(instructions)
		if err != nil {
			log.Panicf("compilation error: %v", err)
		}
	} else {
		// Use provided hex bytecode
		bytecode, err = fvm.HexToBytes(args.Run)
		if err != nil {
			log.Panicf("failed to parse hex bytecode: %v", err)
		}
	}

	if args.Debug {
		fmt.Printf("Compiled bytecode (hex): %s\n", hex.EncodeToString(bytecode))
		disasm, err := fvm.Disassemble(bytecode)
		if err != nil {
			fmt.Printf("Disassembly error: %v\n", err)
		} else {
			fmt.Println("Disassembled bytecode:")
			fmt.Println(disasm)
		}
	}

	// Create VM and program
	prog := fvm.NewProgram(bytecode)
	vm := fvm.New(prog, uint64(args.Gas))

	// Run the VM
	if err := vm.Run(); err != nil {
		log.Panicf("VM execution error: %v", err)
	}

	// Debug output
	if args.Debug {
		fmt.Println("Stack after execution:")
		vm.PrintStackData()
		fmt.Println("Gas remaining:", vm.GasRemaining())

		if args.Caller != "" {
			caller, err := hex.DecodeString(args.Caller)
			if err != nil {
				log.Panicf("invalid caller hex: %v", err)
			}
			fmt.Println("Contract address:")
			vm.PrintContractAddress(caller)
		}
	}

	fmt.Println("Smart contract executed successfully.")
}
