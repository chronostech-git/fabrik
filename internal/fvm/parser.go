package fvm

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	OpCode string
	Arg    string
	Label  string
	Line   int
}

// Parse a .fab contract file.
// See contracts/ folder for example contracts (simple.fab, complex.fab, deposit.fab)
func ParseFile(path string) ([]Instruction, error) {
	if !strings.HasSuffix(path, ".fab") {
		return nil, fmt.Errorf("expected .fab file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ParseLines(lines)
}

func ParseLines(lines []string) ([]Instruction, error) {
	var instructions []Instruction
	labelToInstIndex := make(map[string]int)

	for i, raw := range lines {
		line := cleanLine(raw)
		if line == "" {
			continue
		}

		var label string

		// Only treat as label if colon splits FIRST token
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)

			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			if left != "" && isIdentifier(left) {
				if _, exists := labelToInstIndex[left]; exists {
					return nil, fmt.Errorf("line %d: duplicate label %s", i+1, left)
				}

				labelToInstIndex[left] = len(instructions)
				label = left

				if right == "" {
					instructions = append(instructions, Instruction{
						Label: label,
						Line:  i + 1,
					})
					continue
				}

				line = right
			}
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		op := strings.ToUpper(parts[0])

		if op == "" {
			return nil, fmt.Errorf("line %d: unknown opcode", i+1)
		}

		inst := Instruction{
			OpCode: op,
			Label:  label,
			Line:   i + 1,
		}

		if len(parts) > 1 {
			inst.Arg = parts[1]
		}

		if inst.OpCode == "PUSH" && inst.Arg == "" {
			return nil, fmt.Errorf("line %d: PUSH requires an argument", i+1)
		}

		instructions = append(instructions, inst)
	}

	return instructions, nil
}

func cleanLine(line string) string {
	line = strings.TrimSpace(line)

	// remove BOM if present
	line = strings.TrimPrefix(line, "\ufeff")

	// normalize weird spaces
	line = strings.ReplaceAll(line, "\u00a0", " ")

	for _, sep := range []string{"#", "//", ";"} {
		if idx := strings.Index(line, sep); idx != -1 {
			line = line[:idx]
		}
	}

	return strings.TrimSpace(line)
}

func parseValue(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		return hexToBytes(s[2:])
	}

	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", s)
	}

	return uint64ToBytes(n), nil
}

func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 32)
	for i := 0; i < 8; i++ {
		b[31-i] = byte(v >> (8 * i))
	}
	return b
}

func isIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !(r == '_' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func hexToBytes(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}
