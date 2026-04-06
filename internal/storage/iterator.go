package storage

// Iterator is a Database iterator compatible with both
// memorydb and leveldb.
type Iterator interface {
	Next() bool // Returns true if there is a next item in db
	Key() []byte
	Value() []byte
	Close() error
}
