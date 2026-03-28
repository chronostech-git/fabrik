package fvm

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

func Disassemble(code []byte) (string, error) {
	var out []string

	// First pass: find jump targets
	targets := make(map[int]string)
	labelCount := 0

	pc := 0
	for pc < len(code) {
		op := OpCode(code[pc])
		pc++

		if op == PUSH {
			if pc+32 > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			var v uint256.Int
			v.SetBytes(code[pc : pc+32])

			addr := int(v.Uint64())

			// Heuristic: treat as label if it's inside code
			if addr < len(code) {
				if _, ok := targets[addr]; !ok {
					targets[addr] = fmt.Sprintf("L%d", labelCount)
					labelCount++
				}
			}

			pc += 32
		}
	}

	// Second pass: actual disassembly
	pc = 0
	for pc < len(code) {
		if label, ok := targets[pc]; ok {
			out = append(out, label+":")
		}

		op := OpCode(code[pc])
		pc++

		name := op.String()

		if op == PUSH {
			if pc+32 > len(code) {
				return "", fmt.Errorf("invalid PUSH at pc %d", pc-1)
			}

			var v uint256.Int
			v.SetBytes(code[pc : pc+32])

			val := v.Uint64()

			// Replace with label if known
			if label, ok := targets[int(val)]; ok {
				out = append(out, fmt.Sprintf("%s %s", name, label))
			} else {
				out = append(out, fmt.Sprintf("%s %d", name, val))
			}

			pc += 32
		} else {
			out = append(out, name)
		}
	}

	return strings.Join(out, "\n"), nil
}
