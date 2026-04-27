package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Snapshot represents a captured set of environment variables.
type Snapshot struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// New creates a new Snapshot from a slice of "KEY=VALUE" strings (e.g. os.Environ()).
func New(name string, environ []string) (*Snapshot, error) {
	if name == "" {
		return nil, errors.New("snapshot name must not be empty")
	}

	vars := make(map[string]string, len(environ))
	for _, entry := range environ {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		vars[parts[0]] = parts[1]
	}

	return &Snapshot{
		Name:      name,
		Vars:      vars,
		CreatedAt: time.Now(),
	}, nil
}

// Save writes the snapshot to a JSON file at the given path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write snapshot file: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("snapshot file not found: %s", path)
		}
		return nil, fmt.Errorf("read snapshot file: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &s, nil
}
