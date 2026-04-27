// Package store provides a file-system backed persistence layer for envlock
// snapshots. Snapshots are stored as JSON files inside a configurable
// directory (default: .envlock/ in the working directory).
//
// Typical usage:
//
//	s, err := store.New("") // uses .envlock/ in cwd
//	if err != nil { ... }
//
//	// Persist a snapshot
//	if err := s.Save(snap); err != nil { ... }
//
//	// Retrieve it later
//	loaded, err := s.Load("dev")
//
//	// Enumerate all stored snapshots
//	names, err := s.List()
//
//	// Remove a snapshot
//	if err := s.Delete("dev"); err != nil { ... }
package store
