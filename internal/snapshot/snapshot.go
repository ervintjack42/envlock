package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured set of environment variables.
type Snapshot struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot from the current process environment.
func New(name string, environ []string) (*Snapshot, error) {
	if name == "" {
		return nil, fmt.Errorf("snapshot name must not be empty")
	}

	env := make(map[string]string, len(environ))
	for _, entry := range environ {
		for i := 0; i < len(entry); i++ {
			if entry[i] == '=' {
				env[entry[:i]] = entry[i+1:]
				break
			}
		}
	}

	return &Snapshot{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Env:       env,
	}, nil
}

// Save writes the snapshot to a JSON file at the given path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("encoding snapshot: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening snapshot file: %w", err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("decoding snapshot: %w", err)
	}
	return &s, nil
}
