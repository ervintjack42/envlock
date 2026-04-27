package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/envlock/internal/snapshot"
)

const defaultStoreDir = ".envlock"

// Store manages persisted snapshots on disk.
type Store struct {
	dir string
}

// New creates a Store rooted at the given directory.
// If dir is empty, it defaults to ".envlock" in the current working directory.
func New(dir string) (*Store, error) {
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(cwd, defaultStoreDir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot to the store, keyed by its name.
func (s *Store) Save(snap *snapshot.Snapshot) error {
	if snap.Name == "" {
		return errors.New("snapshot name must not be empty")
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.snapshotPath(snap.Name), data, 0o600)
}

// Load retrieves a snapshot by name from the store.
func (s *Store) Load(name string) (*snapshot.Snapshot, error) {
	if name == "" {
		return nil, errors.New("snapshot name must not be empty")
	}
	data, err := os.ReadFile(s.snapshotPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("snapshot not found: " + name)
		}
		return nil, err
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// List returns the names of all stored snapshots.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a snapshot from the store by name.
func (s *Store) Delete(name string) error {
	if name == "" {
		return errors.New("snapshot name must not be empty")
	}
	err := os.Remove(s.snapshotPath(name))
	if os.IsNotExist(err) {
		return errors.New("snapshot not found: " + name)
	}
	return err
}

func (s *Store) snapshotPath(name string) string {
	return filepath.Join(s.dir, name+".json")
}
