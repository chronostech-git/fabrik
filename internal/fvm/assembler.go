package fvm

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/holiman/uint256"
)

func Assemble(src string) ([]byte, error) {
	lines := strings.Split(src, "\n")
	var bytecode []byte

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		opName := strings.ToUpper(parts[0])

		op, ok := opMap[opName]
		if !ok {
			return nil, fmt.Errorf("line %d: unknown opcode %s", i+1, opName)
		}

		bytecode = append(bytecode, byte(op))

		if op == PUSH {
			if len(parts) != 2 {
				return nil, fmt.Errorf("line %d: PUSH requires an argument", i+1)
			}

			val, err := parseValue(parts[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", i+1, err)
			}

			buf := make([]byte, 32)
			copy(buf[32-len(val):], val)
			bytecode = append(bytecode, buf...)
		}
	}

	return bytecode, nil
}

func parseValue(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		return hexToBytes(s[2:])
	}

	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", s)
	}

	v := uint256.NewInt(0).SetUint64(n)
	return v.Bytes(), nil
}

func hexToBytes(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}
