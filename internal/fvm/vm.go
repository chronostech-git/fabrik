package fvm

import (
	"crypto/sha256"
	"fmt"

	"github.com/chronostech-git/fabrik/internal/state"
	"github.com/holiman/uint256"
)

type VM struct {
	prog     *Program
	pc       int
	stack    *Stack
	gas      uint64
	memory   map[uint64]uint256.Int
	storage  map[uint64]uint256.Int
	dispatch map[OpCode]func(*VM) error
}

func New(prog *Program, state *state.AccountState, gas uint64) *VM {
	vm := &VM{
		prog:     prog,
		pc:       0,
		stack:    NewStack(),
		gas:      gas,
		memory:   make(map[uint64]uint256.Int),
		storage:  make(map[uint64]uint256.Int),
		dispatch: make(map[OpCode]func(*VM) error),
	}
	vm.initDispatch()
	return vm
}

func (vm *VM) initDispatch() {
	vm.dispatch[STOP] = func(vm *VM) error {
		vm.pc = len(vm.prog.code)
		return nil
	}

	vm.dispatch[ADD] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(*uint256.NewInt(0).Add(&a, &b))
		return nil
	}

	vm.dispatch[SUB] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(*uint256.NewInt(0).Sub(&a, &b))
		return nil
	}

	vm.dispatch[MUL] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(*uint256.NewInt(0).Mul(&a, &b))
		return nil
	}

	vm.dispatch[DIV] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(*uint256.NewInt(0).Div(&a, &b))
		return nil
	}

	vm.dispatch[EXP] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(*uint256.NewInt(0).Exp(&a, &b))
		return nil
	}

	// PUSH1–PUSH32
	for i := 0; i < 32; i++ {
		pushOp := 0x12 + i
		pushN := i + 1
		vm.dispatch[OpCode(pushOp)] = func(vm *VM) error {
			if vm.pc+pushN > len(vm.prog.code) {
				return fmt.Errorf("invalid PUSH at pc %d", vm.pc-1)
			}
			val := uint256.NewInt(0)
			val.SetBytes(vm.prog.code[vm.pc : vm.pc+pushN])
			vm.stack.Push(*val)
			vm.pc += pushN
			return nil
		}
	}

	vm.dispatch[POP] = func(vm *VM) error {
		vm.stack.Pop()
		return nil
	}

	vm.dispatch[DUP] = func(vm *VM) error {
		val := vm.stack.Peek()
		vm.stack.Push(*val)
		return nil
	}

	vm.dispatch[SWAP] = func(vm *VM) error {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(a)
		vm.stack.Push(b)
		return nil
	}

	vm.dispatch[JMP] = func(vm *VM) error {
		addr := vm.stack.Pop()
		vm.pc = int(addr.Uint64())
		return nil
	}

	vm.dispatch[JMPI] = func(vm *VM) error {
		addr := vm.stack.Pop()
		cond := vm.stack.Pop()
		if !cond.IsZero() {
			vm.pc = int(addr.Uint64())
		}
		return nil
	}

	vm.dispatch[MSTORE] = func(vm *VM) error {
		addr := vm.stack.Pop()
		val := vm.stack.Pop()
		vm.memory[addr.Uint64()] = val
		return nil
	}

	vm.dispatch[MLOAD] = func(vm *VM) error {
		addr := vm.stack.Pop()
		val, ok := vm.memory[addr.Uint64()]
		if !ok {
			val = *uint256.NewInt(0)
		}
		vm.stack.Push(val)
		return nil
	}

	vm.dispatch[SSTORE] = func(vm *VM) error {
		addr := vm.stack.Pop()
		val := vm.stack.Pop()
		vm.storage[addr.Uint64()] = val
		return nil
	}

	vm.dispatch[SLOAD] = func(vm *VM) error {
		addr := vm.stack.Pop()
		val, ok := vm.storage[addr.Uint64()]
		if !ok {
			val = *uint256.NewInt(0)
		}
		vm.stack.Push(val)
		return nil
	}

	vm.dispatch[SHA256] = func(vm *VM) error {
		val := vm.stack.Pop()
		buf := val.Bytes()
		hash := sha256Sum(buf)
		vm.stack.Push(*uint256.NewInt(0).SetBytes(hash))
		return nil
	}

	vm.dispatch[ADDRESS] = func(vm *VM) error {
		vm.stack.Push(*uint256.NewInt(0))
		return nil
	}

	vm.dispatch[CALLER] = func(vm *VM) error {
		vm.stack.Push(*uint256.NewInt(0))
		return nil
	}
}

// burnGas helper
func (vm *VM) burnGas(amount uint64) error {
	if vm.gas < amount {
		return fmt.Errorf("out of gas")
	}
	vm.gas -= amount
	return nil
}

// sha256 helper
func sha256Sum(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func (vm *VM) Run() error {
	for vm.pc < len(vm.prog.code) {
		op := OpCode(vm.prog.code[vm.pc])
		vm.pc++
		if fn, ok := vm.dispatch[op]; ok {
			if err := vm.burnGas(3); err != nil {
				return err
			}
			if err := fn(vm); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid opcode 0x%x at pc %d", op, vm.pc-1)
		}
	}
	return nil
}

func (vm *VM) PrintStackData() {
	fmt.Println("stack:", vm.stack.data)
}

func (vm *VM) PrintGasRemaining() {
	fmt.Println("gas remaining:", vm.gas)
}
