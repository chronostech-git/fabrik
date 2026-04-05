package memory

import (
	"errors"
	"sync"

	"github.com/chronostech-git/fabrik/internal/storage"
	"github.com/chronostech-git/fabrik/internal/types"
)

var ErrKeyNotFound = errors.New("key not found")
var ErrMemoryDbNotFound = errors.New("memory db was not found")

type MemoryDB struct {
	db   map[string][]byte
	lock sync.RWMutex
}

func New() *MemoryDB {
	return &MemoryDB{
		db: make(map[string][]byte),
	}
}

func (db *MemoryDB) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if v, ok := db.db[string(key)]; ok {
		return types.CopyBytes(v), nil
	}

	return nil, ErrMemoryDbNotFound
}

func (db *MemoryDB) Put(key, value []byte) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	vCopy := types.CopyBytes(value)

	db.db[string(key)] = vCopy

	return nil
}

func (db *MemoryDB) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.db, string(key))
	return nil
}

func (db *MemoryDB) Has(key []byte) (bool, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	_, ok := db.db[string(key)]

	return ok, nil
}

func (db *MemoryDB) Data() map[string][]byte {
	return db.db
}

type memoryIterator struct {
	keys     []string
	db       *MemoryDB
	position int
}

func (db *MemoryDB) NewIterator() storage.Iterator {
	db.lock.RLock()
	defer db.lock.RUnlock()

	keys := make([]string, 0, len(db.Data()))
	for k := range db.Data() {
		keys = append(keys, k)
	}

	return &memoryIterator{
		keys:     keys,
		db:       db,
		position: 0,
	}
}

func (it *memoryIterator) Next() bool {
	return it.position < len(it.keys)
}

func (it *memoryIterator) Key() []byte {
	return []byte(it.keys[it.position])
}

func (it *memoryIterator) Value() []byte {
	it.db.lock.RLock()
	defer it.db.lock.RUnlock()
	data := it.db.Data()
	return data[it.keys[it.position]]
}

// Nothing to clean up for in-memory, therfore it returns "nil"
func (it *memoryIterator) Close() error {
	return nil
}

func (it *memoryIterator) Advance() {
	it.position++
}
