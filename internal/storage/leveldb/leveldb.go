package leveldb

import (
	"sync"

	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type LevelDB struct {
	db   *leveldb.DB
	lock sync.RWMutex
}

func New(file string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{db: db}, nil
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *LevelDB) Put(key, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LevelDB) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LevelDB) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

type levelIterator struct {
	iter iterator.Iterator
}

func (ldb *LevelDB) NewIterator() storage.Iterator {
	iter := ldb.db.NewIterator(nil, nil) // full DB iteration
	return &levelIterator{iter: iter}
}

func (it *levelIterator) Next() bool {
	return it.iter.Next()
}

func (it *levelIterator) Key() []byte {
	return it.iter.Key()
}

func (it *levelIterator) Value() []byte {
	return it.iter.Value()
}

func (it *levelIterator) Close() error {
	it.iter.Release()
	return it.iter.Error()
}
