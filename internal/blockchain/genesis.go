package blockchain

import (
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/types"
)

const GenesisFilename = "genesis.dat"

var ErrGenesisMissing = errors.New("genesis not found")

// ✅ PURE DATA ONLY
type Genesis struct {
	CreationTime int64
	GenesisHash  types.Hash
	Coinbase     types.Address
	InitialValue types.Amount
}

func NewGenesis(time int64, coinbase types.Address, value types.Amount) *Genesis {
	g := &Genesis{
		CreationTime: time,
		Coinbase:     coinbase,
		InitialValue: value,
	}

	g.GenesisHash = g.computeHash()
	return g
}

func (g *Genesis) computeHash() types.Hash {
	enc, _ := rlp.Encode(g)
	return sha256.Sum256(enc)
}

func (g *Genesis) Write(datadir string) error {
	data, _ := rlp.Encode(g)

	if err := os.MkdirAll(datadir, 0700); err != nil {
		return err
	}

	path := filepath.Join(datadir, GenesisFilename)
	return os.WriteFile(path, data, 0600)
}

func LoadGenesis(datadir string) (*Genesis, error) {
	path := filepath.Join(datadir, GenesisFilename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrGenesisMissing
	}

	var g Genesis
	if err := rlp.Decode(data, &g); err != nil {
		return nil, err
	}

	return &g, nil
}
