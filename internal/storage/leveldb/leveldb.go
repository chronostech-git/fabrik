package leveldb

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
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
