package keystore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
)

type FileStore struct {
	dir string
}

func NewFileStore(datadir string) *FileStore {
	return &FileStore{
		dir: filepath.Join(datadir, "keystore"),
	}
}

func (fs *FileStore) StoreKey(k *crypto.Key) error {
	if err := os.MkdirAll(fs.dir, 0700); err != nil {
		return err
	}

	// encode ONLY private key bytes
	privBytes := k.PrivateKey.D.Bytes()

	data, _ := rlp.Encode(privBytes)

	filename := filepath.Join(fs.dir, k.Address.Hex()+".key")
	return os.WriteFile(filename, data, 0600)
}

func (fs *FileStore) GetKey() (*crypto.Key, error) {
	files, err := os.ReadDir(fs.dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("no keys found")
	}

	path := filepath.Join(fs.dir, files[0].Name())
	return loadKey(path)
}

func loadKey(path string) (*crypto.Key, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var privBytes []byte
	if err := rlp.Decode(data, &privBytes); err != nil {
		return nil, err
	}

	priv := crypto.BytesToPrivateKey(privBytes)

	k := &crypto.Key{
		PrivateKey: priv,
	}

	k.Address = crypto.GenerateAddress(&priv.PublicKey)
	return k, nil
}
