package fvm

var gasTable = map[OpCode]uint64{
	STOP: 0,

	ADD: 3,
	SUB: 3,
	MUL: 5,
	DIV: 5,
	EXP: 10,

	PUSH: 3,
	POP:  2,
	DUP:  3,
	SWAP: 3,

	JMP:  8,
	JMPI: 10,

	MLOAD:  5,
	MSTORE: 8,

	SLOAD:  50,
	SSTORE: 20000,

	SHA256: 60,

	ADDRESS: 2,
	CALLER:  2,
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
