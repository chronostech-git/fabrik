package storage

// Basic functions for now.
// This interface WILL be extended to support an iterator as well.
type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	NewIterator() Iterator
}
