package keystore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/crypto"
	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
)

type FileStore struct {
	Dir string
}

// NewFileStore creates a new File Storage given a Data Directory specified using cli/fabkey --datadir <datadir>
func NewFileStore(datadir string) *FileStore {
	return &FileStore{
		Dir: filepath.Join(datadir, "keystore"),
	}
}

// StoreKey is an interface function that writes the key to the .key file
func (fs *FileStore) StoreKey(k *crypto.Key) error {
	if k == nil || k.PrivateKey == nil || k.PrivateKey.D == nil {
		return errors.New("invalid key")
	}

	if err := os.MkdirAll(fs.Dir, 0700); err != nil {
		return err
	}

	privBytes := k.PrivateKey.D.Bytes()
	if len(privBytes) == 0 {
		return errors.New("private key is empty")
	}

	data, err := rlp.Encode(privBytes)
	if err != nil {
		return err
	}

	filename := filepath.Join(fs.Dir, k.Address.Hex()+".key")
	return os.WriteFile(filename, data, 0600)
}

// GetKey loads a key from the file store given the <datadir>/keystore directory exists and has a key.
func (fs *FileStore) GetKey() (*crypto.Key, error) {
	files, err := os.ReadDir(fs.Dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("no keys found")
	}

	path := filepath.Join(fs.Dir, files[0].Name())
	return loadKey(path)
}

// loadKey used to load a key inside GetKey() interface function.
func loadKey(path string) (*crypto.Key, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var privBytes []byte
	if err := rlp.Decode(data, &privBytes); err != nil {
		return nil, err
	}

	if len(privBytes) == 0 {
		return nil, errors.New("decoded private key is empty")
	}

	priv := crypto.BytesToPrivateKey(privBytes)
	k := &crypto.Key{
		PrivateKey: priv,
	}
	k.Address = crypto.GenerateAddress(&priv.PublicKey)
	return k, nil
}
