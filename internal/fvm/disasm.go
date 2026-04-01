package fvm

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

func Disassemble(code []byte) (string, error) {
	var out []string

	targets := make(map[int]string)
	labelCount := 0

	// helper
	getPushSize := func(op OpCode) int {
		if op >= 0x12 && op <= 0x31 {
			return int(op-0x12) + 1
		}
		return 0
	}

	// FIRST PASS (collect jump targets)
	pc := 0
	for pc < len(code) {
		op := OpCode(code[pc])
		pc++

		pushSize := getPushSize(op)

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

	// SECOND PASS (actual disassembly)
	pc = 0
	for pc < len(code) {
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
