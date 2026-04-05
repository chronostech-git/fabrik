package debug

type PrintData interface{}

type Printer interface {
	PrintData()
}
