package fvm

import "github.com/holiman/uint256"

type Stack struct {
	data []uint256.Int
}

func NewStack() *Stack {
	return &Stack{
		data: make([]uint256.Int, 0, 16),
	}
}

func (st *Stack) Push(v uint256.Int) {
	st.data = append(st.data, v)
}

func (st *Stack) Pop() (ret uint256.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) Len() int {
	return len(st.data)
}

func (st *Stack) Peek() *uint256.Int {
	return &st.data[st.Len()-1]
}

func (st *Stack) Back(n int) *uint256.Int {
	return &st.data[st.Len()-n-1]
}
