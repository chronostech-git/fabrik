package main

import (
	"log"

	"github.com/chronostech-git/fabrik/internal/fvm"
)

func main() {
	code, err := fvm.Assemble(`
		PUSH 5
		PUSH 10
		ADD
		PUSH 2
		MUL
		DUP
		PUSH 3
		SUB
		PUSH 0x0a
		EXP
		MSTORE
		MLOAD
		PUSH 0x01
		SSTORE
		SLOAD
		SHA256
		PUSH 0x1234
		ADDRESS
		CALLER
		STOP
	`)
	if err != nil {
		log.Panic(err)
	}

	prog := fvm.NewProgram(code)

	vm := fvm.New(prog, 10000000)

	err = vm.Run()
	if err != nil {
		log.Panic(err)
	}

	vm.PrintGasRemaining()
	vm.PrintStackData()
}
