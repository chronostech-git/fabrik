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

// NewDiskStorage creates an instance of DiskStorage given a
// Storage directory (seperate, or the same as datadir specified in CLI).
func NewDiskStorage(dir string) *DiskStorage {
	return &DiskStorage{
		Directory: dir,
	}
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

// LoadPeers loads the peers in peerJson format from the peers.json file.
func (ds *DiskStorage) LoadPeers() (map[string]peerJson, error) {
	var peers map[string]peerJson

	file := filepath.Join(ds.Directory, peerFilename)

	data, err := os.ReadFile(file)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &peers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal existing peers: %w", err)
		}
	}

	return peers, nil
}
