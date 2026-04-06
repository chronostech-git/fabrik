package fvm

// Program is a wrapper around bytecode.
type Program struct {
	code []byte
}

// NewProgram returns a Program with bytecode.
func NewProgram(code []byte) *Program {
	return &Program{code: code}
}
