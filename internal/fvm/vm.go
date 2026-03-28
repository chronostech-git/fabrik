package fvm

import (
	"crypto/sha256"

	"github.com/chronostech-git/fabrik/internal/types"
	"github.com/holiman/uint256"
)

type FVM struct {
	pc     uint64
	prog   *Program
	stack  *Stack
	memory []byte

	address types.Address
	caller  types.Address

	state StateDB

	stopped bool
}

func New(prog *Program) *FVM {
	return &FVM{
		stack: NewStack(),
		prog:  prog,
	}
}

func (vm *FVM) push(v uint256.Int) {
	vm.stack.Push(v)
}

func (vm *FVM) pop() uint256.Int {
	return vm.stack.Pop()
}

func (vm *FVM) Run() error {
	for !vm.stopped && int(vm.pc) < len(vm.prog.code) {
		op := OpCode(vm.prog.code[vm.pc])
		vm.pc++

		switch op {

		case STOP:
			vm.stopped = true

		case ADD:
			a := vm.pop()
			b := vm.pop()
			var r uint256.Int
			r.Add(&a, &b)
			vm.push(r)

		case SUB:
			a := vm.pop()
			b := vm.pop()
			var r uint256.Int
			r.Sub(&a, &b)
			vm.push(r)

		case MUL:
			a := vm.pop()
			b := vm.pop()
			var r uint256.Int
			r.Mul(&a, &b)
			vm.push(r)

		case DIV:
			a := vm.pop()
			b := vm.pop()
			var r uint256.Int
			if b.IsZero() {
				r.Clear()
			} else {
				r.Div(&a, &b)
			}
			vm.push(r)

		case EXP:
			a := vm.pop()
			b := vm.pop()
			var r uint256.Int
			r.Exp(&a, &b)
			vm.push(r)

		case PUSH:
			if int(vm.pc+32) > len(vm.prog.code) {
				return ErrOutOfBounds
			}
			var val uint256.Int
			val.SetBytes(vm.prog.code[vm.pc : vm.pc+32])
			vm.pc += 32
			vm.push(val)

		case POP:
			vm.pop()

		case DUP:
			v := *vm.stack.Peek()
			vm.push(v)

		case SWAP:
			a := vm.pop()
			b := vm.pop()
			vm.push(a)
			vm.push(b)

		case JMP:
			dest := vm.pop()
			vm.pc = dest.Uint64()

		case JMPI:
			dest := vm.pop()
			cond := vm.pop()
			if !cond.IsZero() {
				vm.pc = dest.Uint64()
			}

		case MSTORE:
			v := vm.pop()
			offset := v.Uint64()
			value := vm.pop()
			if int(offset+32) > len(vm.memory) {
				newMem := make([]byte, offset+32)
				copy(newMem, vm.memory)
				vm.memory = newMem
			}
			copy(vm.memory[offset:], value.Bytes())

		case MLOAD:
			v := vm.pop()
			offset := v.Uint64()
			var val uint256.Int
			if int(offset+32) <= len(vm.memory) {
				val.SetBytes(vm.memory[offset : offset+32])
			}
			vm.push(val)

		case SSTORE:
			key := vm.pop()
			val := vm.pop()
			vm.state.SetState(vm.address, key, val)

		case SLOAD:
			key := vm.pop()
			val := vm.state.GetState(vm.address, key)
			vm.push(val)

		case SHA256:
			val := vm.pop()
			hash := sha256.Sum256(val.Bytes())
			var out uint256.Int
			out.SetBytes(hash[:])
			vm.push(out)

		case ADDRESS:
			var v uint256.Int
			v.SetBytes(vm.address.Bytes())
			vm.push(v)

		case CALLER:
			var v uint256.Int
			v.SetBytes(vm.caller.Bytes())
			vm.push(v)

		default:
			return ErrInvalidOpcode
		}
	}
	return nil
}
