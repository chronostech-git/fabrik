package fvm

import (
	"fmt"
	"strconv"

	"github.com/holiman/uint256"
)

func Compile(instructions []Instruction) ([]byte, error) {
	var bytecode []byte

	instToPC := make([]int, len(instructions))
	pc := 0

	for i, inst := range instructions {
		instToPC[i] = pc

		if inst.OpCode == "" {
			continue
		}

		op, ok := opMap[inst.OpCode]
		if !ok {
			return nil, fmt.Errorf("line %d: unknown opcode %s", inst.Line, inst.OpCode)
		}

		if op == PUSH {
			pc += 33
		} else {
			pc += 1
		}
	}

	labelToPC := make(map[string]int)
	for i, inst := range instructions {
		if inst.Label != "" {
			labelToPC[inst.Label] = instToPC[i]
		}
	}

	for _, inst := range instructions {
		if inst.OpCode == "" {
			continue
		}

		op, ok := opMap[inst.OpCode]
		if !ok {
			return nil, fmt.Errorf("line %d: unknown opcode %s", inst.Line, inst.OpCode)
		}

		bytecode = append(bytecode, byte(op))

		if op == PUSH {
			if inst.Arg == "" {
				return nil, fmt.Errorf("line %d: PUSH requires argument", inst.Line)
			}

			if addr, ok := labelToPC[inst.Arg]; ok {
				var v uint256.Int
				v.SetUint64(uint64(addr))
				b := v.Bytes()
				buf := make([]byte, 32)
				copy(buf[32-len(b):], b)
				bytecode = append(bytecode, buf...)
				continue
			}

			val, err := parseValue(inst.Arg)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", inst.Line, err)
			}

			buf := make([]byte, 32)
			copy(buf[32-len(val):], val)
			bytecode = append(bytecode, buf...)
		}
	}

	return bytecode, nil
}

func isLabel(s string) bool {
	if _, err := strconv.ParseUint(s, 10, 64); err == nil {
		return false
	}
	return true
}
