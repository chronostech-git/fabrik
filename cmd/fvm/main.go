package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/chronostech-git/fabrik/internal/fvm"
)

func main() {
	code, err := fvm.Assemble(`
		PUSH 5
		PUSH 10
		ADD
		STOP
	`)
	if err != nil {
		log.Panic(err)
	}

	prog := fvm.NewProgram(code)

	vm := fvm.New(prog)

	err = vm.Run()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(hex.EncodeToString(code))
}
