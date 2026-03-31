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

// Genesis holds all the nessecary fields to kickstart the blockchain.
// The most important fields being coinbase address (Coinbase): This is the address of the wallet key used to create the coinbase transaction. initial value (InitialValue): This is to determine how many FAB coins are created and put into circulation on genesis creation.
type Genesis struct {
	CreationTime int64
	GenesisHash  types.Hash
	Coinbase     types.Address
	Txs          []*Transaction
	InitialValue types.Amount
}

// Create a new genesis given time of creation, coinbase wallet key address, and the initial value.
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

// Write the genesis to the disk.
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

// Load genesis block from the disk.
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
