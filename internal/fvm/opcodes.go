package fvm

import "fmt"

type OpCode byte

// Arithmatic
const (
	STOP OpCode = 0x0
	ADD  OpCode = 0x1
	SUB  OpCode = 0x2
	MUL  OpCode = 0x3
	DIV  OpCode = 0x4
	EXP  OpCode = 0x5
)

// Stack operations
const (
	POP  OpCode = 0x13
	DUP  OpCode = 0x14
	SWAP OpCode = 0x15
	JMP  OpCode = 0x16
	JMPI OpCode = 0x17
)

// PUSH1 ... PUSH32
const (
	PUSH1  OpCode = 0x60
	PUSH2  OpCode = 0x61
	PUSH3  OpCode = 0x62
	PUSH4  OpCode = 0x63
	PUSH5  OpCode = 0x64
	PUSH6  OpCode = 0x65
	PUSH7  OpCode = 0x66
	PUSH8  OpCode = 0x67
	PUSH9  OpCode = 0x68
	PUSH10 OpCode = 0x69
	PUSH11 OpCode = 0x6a
	PUSH12 OpCode = 0x6b
	PUSH13 OpCode = 0x6c
	PUSH14 OpCode = 0x6d
	PUSH15 OpCode = 0x6e
	PUSH16 OpCode = 0x6f
	PUSH17 OpCode = 0x70
	PUSH18 OpCode = 0x71
	PUSH19 OpCode = 0x72
	PUSH20 OpCode = 0x73
	PUSH21 OpCode = 0x74
	PUSH22 OpCode = 0x75
	PUSH23 OpCode = 0x76
	PUSH24 OpCode = 0x77
	PUSH25 OpCode = 0x78
	PUSH26 OpCode = 0x79
	PUSH27 OpCode = 0x7a
	PUSH28 OpCode = 0x7b
	PUSH29 OpCode = 0x7c
	PUSH30 OpCode = 0x7d
	PUSH31 OpCode = 0x7e
	PUSH32 OpCode = 0x7f
)

// Environment / system
const (
	ADDRESS        OpCode = 0x30
	BALANCE        OpCode = 0x31
	ORIGIN         OpCode = 0x32
	CALLER         OpCode = 0x33
	CALLVALUE      OpCode = 0x34
	CALLDATALOAD   OpCode = 0x35
	CALLDATASIZE   OpCode = 0x36
	CALLDATACOPY   OpCode = 0x37
	CODESIZE       OpCode = 0x38
	CODECOPY       OpCode = 0x39
	GASPRICE       OpCode = 0x3a
	EXTCODESIZE    OpCode = 0x3b
	EXTCODECOPY    OpCode = 0x3c
	RETURNDATASIZE OpCode = 0x3d
	RETURNDATACOPY OpCode = 0x3e
	EXTCODEHASH    OpCode = 0x3f
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

func (op OpCode) String() string {
	switch op {
	// Arithmetic
	case STOP:
		return "STOP"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case EXP:
		return "EXP"

	// Stack
	case POP:
		return "POP"
	case DUP:
		return "DUP"
	case SWAP:
		return "SWAP"

	// PUSH1..PUSH32
	case PUSH1:
		return "PUSH1"
	case PUSH2:
		return "PUSH2"
	case PUSH3:
		return "PUSH3"
	case PUSH4:
		return "PUSH4"
	case PUSH5:
		return "PUSH5"
	case PUSH6:
		return "PUSH6"
	case PUSH7:
		return "PUSH7"
	case PUSH8:
		return "PUSH8"
	case PUSH9:
		return "PUSH9"
	case PUSH10:
		return "PUSH10"
	case PUSH11:
		return "PUSH11"
	case PUSH12:
		return "PUSH12"
	case PUSH13:
		return "PUSH13"
	case PUSH14:
		return "PUSH14"
	case PUSH15:
		return "PUSH15"
	case PUSH16:
		return "PUSH16"
	case PUSH17:
		return "PUSH17"
	case PUSH18:
		return "PUSH18"
	case PUSH19:
		return "PUSH19"
	case PUSH20:
		return "PUSH20"
	case PUSH21:
		return "PUSH21"
	case PUSH22:
		return "PUSH22"
	case PUSH23:
		return "PUSH23"
	case PUSH24:
		return "PUSH24"
	case PUSH25:
		return "PUSH25"
	case PUSH26:
		return "PUSH26"
	case PUSH27:
		return "PUSH27"
	case PUSH28:
		return "PUSH28"
	case PUSH29:
		return "PUSH29"
	case PUSH30:
		return "PUSH30"
	case PUSH31:
		return "PUSH31"
	case PUSH32:
		return "PUSH32"

	// Control flow
	case JMP:
		return "JMP"
	case JMPI:
		return "JMPI"

	// Memory / storage
	case MSTORE:
		return "MSTORE"
	case MLOAD:
		return "MLOAD"
	case SSTORE:
		return "SSTORE"
	case SLOAD:
		return "SLOAD"

	// Hashing / Crypto
	case SHA256:
		return "SHA256"

	// Environment / System
	case ADDRESS:
		return "ADDRESS"
	case BALANCE:
		return "BALANCE"
	case ORIGIN:
		return "ORIGIN"
	case CALLER:
		return "CALLER"
	case CALLVALUE:
		return "CALLVALUE"
	case CALLDATALOAD:
		return "CALLDATALOAD"
	case CALLDATASIZE:
		return "CALLDATASIZE"
	case CALLDATACOPY:
		return "CALLDATACOPY"
	case CODESIZE:
		return "CODESIZE"
	case CODECOPY:
		return "CODECOPY"
	case GASPRICE:
		return "GASPRICE"
	case EXTCODESIZE:
		return "EXTCODESIZE"
	case EXTCODECOPY:
		return "EXTCODECOPY"
	case RETURNDATASIZE:
		return "RETURNDATASIZE"
	case RETURNDATACOPY:
		return "RETURNDATACOPY"
	case EXTCODEHASH:
		return "EXTCODEHASH"

	default:
		return fmt.Sprintf("UNKNOWN(0x%x)", byte(op))
	}
}
