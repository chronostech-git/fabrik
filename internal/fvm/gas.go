package fvm

// gasTable is a map consisting of an OpCode,
// and it's gas cost.
var gasTable = map[OpCode]uint64{
	STOP:           0,
	ADD:            3,
	SUB:            3,
	MUL:            5,
	DIV:            5,
	EXP:            10,
	POP:            2,
	DUP:            3,
	SWAP:           3,
	PUSH1:          3,
	PUSH2:          3,
	PUSH3:          3,
	PUSH4:          3,
	PUSH5:          3,
	PUSH6:          3,
	PUSH7:          3,
	PUSH8:          3,
	PUSH9:          3,
	PUSH10:         3,
	PUSH11:         3,
	PUSH12:         3,
	PUSH13:         3,
	PUSH14:         3,
	PUSH15:         3,
	PUSH16:         3,
	PUSH17:         3,
	PUSH18:         3,
	PUSH19:         3,
	PUSH20:         3,
	PUSH21:         3,
	PUSH22:         3,
	PUSH23:         3,
	PUSH24:         3,
	PUSH25:         3,
	PUSH26:         3,
	PUSH27:         3,
	PUSH28:         3,
	PUSH29:         3,
	PUSH30:         3,
	PUSH31:         3,
	PUSH32:         3,
	JMP:            8,
	JMPI:           10,
	MLOAD:          3,
	MSTORE:         3,
	SLOAD:          50,
	SSTORE:         20000,
	SHA256:         60,
	ADDRESS:        2,
	BALANCE:        400,
	ORIGIN:         2,
	CALLER:         2,
	CALLVALUE:      2,
	CALLDATALOAD:   3,
	CALLDATASIZE:   2,
	CALLDATACOPY:   3,
	CODESIZE:       2,
	CODECOPY:       3,
	GASPRICE:       2,
	EXTCODESIZE:    700,
	EXTCODECOPY:    700,
	RETURNDATASIZE: 2,
	RETURNDATACOPY: 3,
	EXTCODEHASH:    400,
}

func burnGas(vm *VM, op OpCode) {
	cost, ok := gasTable[op]
	if !ok {
		cost = 1
	}
	if vm.gas < cost {
		panic("out of gas")
	}
	vm.gas -= cost
}
