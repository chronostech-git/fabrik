package fvm

type Program struct {
	code []byte
}

func NewProgram(code []byte) *Program {
	return &Program{code: code}
}
