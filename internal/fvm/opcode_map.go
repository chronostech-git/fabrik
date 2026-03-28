package fvm

var opMap = map[string]OpCode{
	"STOP":    STOP,
	"ADD":     ADD,
	"SUB":     SUB,
	"MUL":     MUL,
	"DIV":     DIV,
	"EXP":     EXP,
	"PUSH":    PUSH,
	"POP":     POP,
	"DUP":     DUP,
	"SWAP":    SWAP,
	"JMP":     JMP,
	"JMPI":    JMPI,
	"MSTORE":  MSTORE,
	"MLOAD":   MLOAD,
	"SSTORE":  SSTORE,
	"SLOAD":   SLOAD,
	"SHA256":  SHA256,
	"ADDRESS": ADDRESS,
	"CALLER":  CALLER,
}
