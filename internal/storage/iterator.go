package storage

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Close() error
}
