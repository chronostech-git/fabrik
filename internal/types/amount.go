package types

import (
	"errors"
	"math/big"
)

type Amount struct {
	i *big.Int
}

var (
	ErrNegativeAmount = errors.New("amount cannot be negative")
)

// NewAmount creates an instance of Amount given a value of int64
func NewAmount(x int64) Amount {
	if x < 0 {
		panic(ErrNegativeAmount)
	}
	return Amount{i: big.NewInt(x)}
}

// NewAmountFromBig creates an instance of Amount given a value of *big.Int
func NewAmountFromBig(x *big.Int) Amount {
	if x == nil {
		return Amount{i: new(big.Int)}
	}
	if x.Sign() < 0 {
		panic(ErrNegativeAmount)
	}
	return Amount{i: new(big.Int).Set(x)}
}

// ZeroAmount returns an empty amount (value of 0)
func ZeroAmount() Amount {
	return Amount{i: new(big.Int)}
}

func (a Amount) Big() *big.Int {
	if a.i == nil {
		return new(big.Int)
	}
	return new(big.Int).Set(a.i)
}

func (a Amount) String() string {
	if a.i == nil {
		return "0"
	}
	return a.i.String()
}

func (a Amount) IsZero() bool {
	return a.i == nil || a.i.Sign() == 0
}

func (a Amount) LessThan(b Amount) bool {
	return a.i.Cmp(b.i) < 0
}

func (a Amount) GreaterThan(b Amount) bool {
	return a.i.Cmp(b.i) > 0
}

func (a Amount) Cmp(b Amount) int {
	return a.Big().Cmp(b.Big())
}

func (a Amount) Equal(b Amount) bool {
	return a.Cmp(b) == 0
}

func (a Amount) Add(b Amount) Amount {
	return Amount{i: new(big.Int).Add(a.Big(), b.Big())}
}

func (a Amount) Sub(b Amount) (Amount, error) {
	if a.Cmp(b) < 0 {
		return Amount{}, errors.New("insufficient amount")
	}
	return Amount{i: new(big.Int).Sub(a.Big(), b.Big())}, nil
}

func (a Amount) Mul(b Amount) Amount {
	return Amount{i: new(big.Int).Mul(a.Big(), b.Big())}
}

func (a Amount) Div(b Amount) (Amount, error) {
	if b.IsZero() {
		return Amount{}, errors.New("division by zero")
	}
	return Amount{i: new(big.Int).Div(a.Big(), b.Big())}, nil
}

func (a Amount) Bytes() []byte {
	return a.Big().Bytes()
}

func BytesToAmount(b []byte) Amount {
	if len(b) == 0 {
		return ZeroAmount()
	}
	return Amount{i: new(big.Int).SetBytes(b)}
}
