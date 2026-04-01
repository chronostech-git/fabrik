package fvm

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

// Takes in bytecode and converts it into a set of instructions.
// NOTE: This will be useful when a contract can be executed only given the bytecode.
// The following example using EVM (Ethereum Virtual Machine), would PUSH 2 and 3, add them,
// and then stop the program (a similar feature will be used for Fabrik Virtual Machine).
// EVM example: ./evm run --debug 60026003016000
// FVM example: cli/fvm --run 60026003016000 --debug
func Disassemble(code []byte) (string, error) {
	var out []string

	targets := make(map[int]string)
	labelCount := 0

	// FIRST PASS
	pc := 0
	for pc < len(code) {
		op := OpCode(code[pc])
		pc++

		pushSize := 0
		if op == PUSH {
			pushSize = 32
		}

		if pushSize > 0 {
			if pc+pushSize > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			var v uint256.Int
			v.SetBytes(code[pc : pc+pushSize])

			addr := int(v.Uint64())
			if addr < len(code) {
				if _, ok := targets[addr]; !ok {
					targets[addr] = fmt.Sprintf("L%d", labelCount)
					labelCount++
				}
			}

			pc += pushSize
		}
	}

	// SECOND PASS
	pc = 0
	for pc < len(code) {
		if label, ok := targets[pc]; ok {
			out = append(out, label+":")
		}

		op := OpCode(code[pc])
		pc++

		name := op.String()
		pushSize := 0
		if op == PUSH {
			pushSize = 32
		}

		if pushSize > 0 {
			if pc+pushSize > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			var v uint256.Int
			v.SetBytes(code[pc : pc+pushSize])
			val := v.Uint64()

			if label, ok := targets[int(val)]; ok {
				out = append(out, fmt.Sprintf("%s %s", name, label))
			} else {
				out = append(out, fmt.Sprintf("%s %d", name, val))
			}

			pc += pushSize
		} else {
			out = append(out, name)
		}
	}

	return strings.Join(out, "\n"), nil
}
