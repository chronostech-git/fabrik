package fvm

import (
	"log"

	"github.com/holiman/uint256"
)

type opcodeFunc func(vm *VM)

var DispatchTable map[OpCode]opcodeFunc

func init() {
	DispatchTable = map[OpCode]opcodeFunc{
		STOP:    opStop,
		ADD:     opAdd,
		SUB:     opSub,
		MUL:     opMul,
		DIV:     opDiv,
		EXP:     opExp,
		PUSH:    opPush,
		POP:     opPop,
		DUP:     opDup,
		SWAP:    opSwap,
		JMP:     opJmp,
		JMPI:    opJmpi,
		MSTORE:  opMstore,
		SSTORE:  opSstore,
		MLOAD:   opMload,
		SLOAD:   opSload,
		SHA256:  opSha256,
		ADDRESS: opAddress,
		CALLER:  opCaller,
	}
}

func opStop(vm *VM) {
	vm.pc = len(vm.prog.code)
}

func opAdd(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	res := uint256.NewInt(0).Add(&a, &b)
	vm.stack.Push(*res)
	vm.burnGas(3)
}

func opSub(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	res := uint256.NewInt(0).Sub(&a, &b)
	vm.stack.Push(*res)
	vm.burnGas(3)
}

func opMul(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	res := uint256.NewInt(0).Mul(&a, &b)
	vm.stack.Push(*res)
	vm.burnGas(5)
}

func opDiv(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	res := uint256.NewInt(0).Div(&a, &b)
	vm.stack.Push(*res)
	vm.burnGas(5)
}

func opExp(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	res := uint256.NewInt(0).Exp(&a, &b)
	vm.stack.Push(*res)
	vm.burnGas(10)
}

func opPush(vm *VM) {
	if vm.pc+32 >= len(vm.prog.code) {
		log.Panic("invalid PUSH, not enough bytes")
	}
	val := vm.prog.code[vm.pc+1 : vm.pc+33]
	vm.pc += 32
	i := uint256.NewInt(0).SetBytes(val)
	vm.stack.Push(*i)
	vm.burnGas(3)
}

func opPop(vm *VM) {
	vm.stack.Pop()
	vm.burnGas(2)
}

func opDup(vm *VM) {
	val := vm.stack.Back(1)
	vm.stack.Push(*val)
	vm.burnGas(3)
}

func opSwap(vm *VM) {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	vm.stack.Push(a)
	vm.stack.Push(b)
	vm.burnGas(3)
}

func opJmp(vm *VM) {
	addr := vm.stack.Pop()
	vm.pc = int(addr.Uint64())
	vm.burnGas(8)
}

func opJmpi(vm *VM) {
	addr := vm.stack.Pop()
	cond := vm.stack.Pop()
	if !cond.IsZero() {
		vm.pc = int(addr.Uint64())
	}
	vm.burnGas(10)
}

func opMstore(vm *VM) {
	addr := vm.stack.Pop()
	val := vm.stack.Pop()
	tmp := vm.memory[addr.Uint64()]
	tmpPtr := &tmp
	tmpPtr.Add(tmpPtr, &val)
	vm.memory[addr.Uint64()] = *tmpPtr
	vm.burnGas(5)
}

func opMload(vm *VM) {
	addr := vm.stack.Pop()
	val := vm.memory[addr.Uint64()]
	vm.stack.Push(val)
	vm.burnGas(5)
}

func opSstore(vm *VM) {
	addr := vm.stack.Pop()
	val := vm.stack.Pop()
	tmp := vm.storage[addr.Uint64()]
	tmpPtr := &tmp
	tmpPtr.Add(tmpPtr, &val)
	vm.storage[addr.Uint64()] = *tmpPtr
	vm.burnGas(20)
}

func opSload(vm *VM) {
	addr := vm.stack.Pop()
	val := vm.storage[addr.Uint64()]
	vm.stack.Push(val)
	vm.burnGas(20)
}

func opSha256(vm *VM) {
	val := vm.stack.Pop()
	// hash is currently placeholder
	vm.stack.Push(val)
	vm.burnGas(30)
}

func opAddress(vm *VM) {
	vm.stack.Push(*uint256.NewInt(0))
	vm.burnGas(2)
}

func opCaller(vm *VM) {
	vm.stack.Push(*uint256.NewInt(0))
	vm.burnGas(2)
}
