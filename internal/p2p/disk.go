package p2p

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	peerFilename = "peers.json"
)

type DiskStorage struct {
	Directory string
}

func NewDiskStorage(dir string) *DiskStorage {
	return &DiskStorage{
		Directory: dir,
	}
}

func (ds *DiskStorage) createPeerFile() error {
	file := filepath.Join(ds.Directory, peerFilename) // <datadir>/peers.json

	_, err := os.Create(file)
	if err != nil {
		return err
	}

	return nil
}

// WritePeer
// If a node connects to the network that is not already saved
// to peers.json, it will write the peer to the disk at <datadir>/peers.json
func (ds *DiskStorage) WritePeer(p *Peer) error {
	file := filepath.Join(ds.Directory, peerFilename)

	var peers []peerJson

	data, err := os.ReadFile(file)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &peers); err != nil {
			return fmt.Errorf("failed to unmarshal existing peers: %w", err)
		}
	}

	peers = append(peers, PeerToJson(p))

	newData, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(file, newData, 0644); err != nil {
		return err
	}

	return nil
}
