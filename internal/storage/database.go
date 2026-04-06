package storage

// Database interface specifies db functions and iterator functions for
// memorydb and leveldb
type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	NewIterator() Iterator
}
