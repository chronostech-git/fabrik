package fvm

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

// Disassemble converts bytecode into human-readable stack-based format.
func Disassemble(code []byte) (string, error) {
	var out []string

	// Map for labels (jump targets)
	targets := make(map[int]string)
	labelCount := 0

	// Helper: returns size of data for PUSH1..PUSH32
	getPushSize := func(op OpCode) int {
		if op >= 0x60 && op <= 0x7f {
			return int(op - 0x60 + 1)
		}
		return 0
	}

	// FIRST PASS: identify jump targets
	pc := 0
	for pc < len(code) {
		op := OpCode(code[pc])
		pc++

		pushSize := getPushSize(op)
		if pushSize > 0 {
			if pc+pushSize > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			// Read value to check if it is a jump target
			var val uint256.Int
			val.SetBytes(code[pc : pc+pushSize])
			addr := int(val.Uint64())

			if addr < len(code) {
				if _, ok := targets[addr]; !ok {
					targets[addr] = fmt.Sprintf("L%d", labelCount)
					labelCount++
				}
			}

			pc += pushSize
		}
	}

	// SECOND PASS: disassemble
	pc = 0
	for pc < len(code) {
		// Print label if current PC is a jump target
		if label, ok := targets[pc]; ok {
			out = append(out, label+":")
		}

		op := OpCode(code[pc])
		pc++

		name := op.String()
		pushSize := getPushSize(op)

		if pushSize > 0 {
			if pc+pushSize > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			var val uint256.Int
			val.SetBytes(code[pc : pc+pushSize])
			value := val.Uint64()

			// If value points to a label, use it
			if label, ok := targets[int(value)]; ok {
				out = append(out, fmt.Sprintf("%s %s", name, label))
			} else {
				out = append(out, fmt.Sprintf("%s %d", name, value))
			}

			pc += pushSize
		} else {
			out = append(out, name)
		}
	}

	return strings.Join(out, "\n"), nil
}
