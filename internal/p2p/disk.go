package p2p

import (
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

func (ds *DiskStorage) WritePeer(p *Peer) error {
	file := filepath.Join(ds.Directory, peerFilename)

	// Ensure file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = ds.createPeerFile()
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data := PeerToJson(p)

	n, err := f.WriteString(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf(
			"Failed to write peer to disk %s. Expected data length of %d, got %d",
			file, len(data), n,
		)
	}

	return nil
}
