package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chronostech-git/fabrik/internal/serialize/rlp"
	"github.com/chronostech-git/fabrik/internal/types"
)

const GenesisFilename = "genesis.dat"

var ErrGenesisMissing = errors.New("genesis not found")

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
	data, err := rlp.Encode(g)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(datadir, 0700); err != nil {
		return err
	}

	path := filepath.Join(datadir, GenesisFilename)
	return os.WriteFile(path, data, 0600)
}

func LoadGenesis(root string) (*Genesis, error) {
	path := filepath.Join(root, "genesis", GenesisFilename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrGenesisMissing
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("genesis.dat is empty")
	}

	var g Genesis
	if err := rlp.Decode(data, &g); err != nil {
		return nil, fmt.Errorf("rlp decode failed: %w", err)
	}

	return &g, nil
}
