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

// Genesis holds all the necessary fields to kickstart the blockchain.
// The most important fields being Coinbase address: this is the address of the fabkey key
// used to create the Coinbase transaction. initial value (InitialValue): This is to determine how many
// FAB coins are created and put into circulation on genesis creation.
type Genesis struct {
	CreationTime int64
	GenesisHash  types.Hash
	Coinbase     types.Address
	Txs          []*Transaction
	InitialValue types.Amount
}

// NewGenesis creates the Genesis block and sets the hash of Genesis block.
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

// ToBlock As mentioned in chain.go, Genesis is a block BUT a Block is not a genesis.
// Converts Genesis to a traditional blockchain.Block.
func (g *Genesis) ToBlock() (*Block, error) {
	enc, err := rlp.Encode(g.Txs)
	if err != nil {
		return nil, err
	}
	txRoot := sha256.Sum256(enc)

	genesisHeader := BlockHeader{
		PrevHash:  types.Empty32(),
		Timestamp: g.CreationTime,
		TxRoot:    txRoot,
		Height:    0,
		GasLimit:  0, // TODO this should be set by the user upon --new, or chain creation
		GasUsed:   0,
		BaseFee:   types.ZeroAmount(),
	}

	return &Block{
		Header:    genesisHeader,
		Txs:       g.Txs,
		StateRoot: types.Empty32(),
		Hash:      g.GenesisHash,
	}, nil
}

// Write does the same operation as Write in chain.go.
// Write in this context (genesis.go), is a compatible version
// of Write for the Genesis struct.
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

// LoadGenesis loads the genesis from the disk (.dat file)
// if present. If not present, or the file is empty for some reason,
// it will return an error.
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
