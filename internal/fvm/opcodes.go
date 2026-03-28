package fvm

type OpCode byte

// Arithmatic
const (
	STOP OpCode = 0x0
	ADD  OpCode = 0x1
	SUB  OpCode = 0x2
	MUL  OpCode = 0x3
	DIV  OpCode = 0x4
	EXP  OpCode = 0x5
	// ends at 0xb (+2 after 0x9)
)

// Stack operations
const (
	PUSH OpCode = 0x12
	POP  OpCode = 0x13
	DUP  OpCode = 0x14
	SWAP OpCode = 0x15
)

// Control flow
const (
	JMP  OpCode = 0x21
	JMPI OpCode = 0x22
)

// Memory / storage
const (
	MSTORE OpCode = 0x41
	SSTORE OpCode = 0x42
	MLOAD  OpCode = 0x43
	SLOAD  OpCode = 0x44
)

// Hashing / Crypto
const (
	SHA256 OpCode = 0x55
)

const (
	ADDRESS OpCode = 0x65
	CALLER  OpCode = 0x66
)
