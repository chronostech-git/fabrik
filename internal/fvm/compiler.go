package fvm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/holiman/uint256"
)

// Compile function takes in the broken down instructions of a .fab contract file
// and returns the bytecode which is then ran using vm.Run().
// FVM creation -> parser -> compiler -> vm.Run().
func Compile(instructions []Instruction) ([]byte, error) {
	var bytecode []byte

	instToPC := make([]int, len(instructions))
	pc := 0

	for i, inst := range instructions {
		instToPC[i] = pc

		if inst.OpCode == "" {
			continue
		}

		_, ok := opMap[inst.OpCode]
		if !ok {
			return nil, fmt.Errorf("line %d: unknown opcode %s", inst.Line, inst.OpCode)
		}

		if strings.HasPrefix(inst.OpCode, "PUSH") {
			// Determine N for PUSHn (1..32)
			nStr := strings.TrimPrefix(inst.OpCode, "PUSH")
			n, err := strconv.Atoi(nStr)
			if err != nil || n < 1 || n > 32 {
				return nil, fmt.Errorf("line %d: invalid PUSH opcode %s", inst.Line, inst.OpCode)
			}
			pc += 1 + n // 1 byte for opcode + n bytes for data
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

		if strings.HasPrefix(inst.OpCode, "PUSH") {
			nStr := strings.TrimPrefix(inst.OpCode, "PUSH")
			n, _ := strconv.Atoi(nStr) // already validated above

			if inst.Arg == "" {
				return nil, fmt.Errorf("line %d: %s requires argument", inst.Line, inst.OpCode)
			}

			var buf []byte
			if addr, ok := labelToPC[inst.Arg]; ok {
				var v uint256.Int
				v.SetUint64(uint64(addr))
				b := v.Bytes()
				buf = make([]byte, n)
				copy(buf[n-len(b):], b)
			} else {
				val, err := parseValue(inst.Arg)
				if err != nil {
					return nil, fmt.Errorf("line %d: %v", inst.Line, err)
				}
				buf = make([]byte, n)
				copy(buf[n-len(val):], val)
			}

			bytecode = append(bytecode, buf...)
		}
	}

	return bytecode, nil
}
