package fvm

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

func Compile(instructions []Instruction) ([]byte, error) {
	var bytecode []byte

	instToPC := make([]int, len(instructions))
	pc := 0

	// FIRST PASS: calculate program counter offsets
	for i, inst := range instructions {
		instToPC[i] = pc

		if inst.OpCode == "" {
			continue
		}

		if strings.HasPrefix(inst.OpCode, "PUSH") {
			// Temporarily assume max 32 bytes
			pc += 1 + 32
		} else {
			pc += 1
		}
	}

	// Map labels to PC offsets
	labelToPC := make(map[string]int)
	for i, inst := range instructions {
		if inst.Label != "" {
			labelToPC[inst.Label] = instToPC[i]
		}
	}

	// SECOND PASS: generate bytecode
	for _, inst := range instructions {
		if inst.OpCode == "" {
			continue
		}

		if strings.HasPrefix(inst.OpCode, "PUSH") {
			var valBytes []byte

			// Resolve argument
			if addr, ok := labelToPC[inst.Arg]; ok {
				var v uint256.Int
				v.SetUint64(uint64(addr))
				valBytes = v.Bytes()
			} else {
				v, err := parseValue(inst.Arg)
				if err != nil {
					return nil, fmt.Errorf("line %d: %v", inst.Line, err)
				}
				valBytes = v
			}

			// Find smallest PUSHn
			n := len(valBytes)
			for n > 1 && valBytes[0] == 0 {
				valBytes = valBytes[1:]
				n--
			}
			if n < 1 || n > 32 {
				return nil, fmt.Errorf("line %d: value too large for PUSHn", inst.Line)
			}

			op := OpCode(0x12 + n - 1) // PUSH1 = 0x12
			bytecode = append(bytecode, byte(op))

			buf := make([]byte, n)
			copy(buf[n-len(valBytes):], valBytes)
			bytecode = append(bytecode, buf...)

		} else {
			op, ok := opMap[inst.OpCode]
			if !ok {
				return nil, fmt.Errorf("line %d: unknown opcode %s", inst.Line, inst.OpCode)
			}
			bytecode = append(bytecode, byte(op))
		}
	}

	return bytecode, nil
}
