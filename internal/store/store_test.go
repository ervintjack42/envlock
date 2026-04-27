package store_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envlock/internal/snapshot"
	"github.com/envlock/internal/store"
)

func newTempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(filepath.Join(dir, "envlock"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func makeSnap(name string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Vars:      map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	s := newTempStore(t)
	snap := makeSnap("dev")

	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != snap.Name {
		t.Errorf("name: got %q, want %q", got.Name, snap.Name)
	}
	if got.Vars["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", got.Vars["FOO"], "bar")
	}
}

func TestStore_LoadMissing(t *testing.T) {
	s := newTempStore(t)
	_, err := s.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestStore_SaveEmptyNameError(t *testing.T) {
	s := newTempStore(t)
	snap := makeSnap("")
	if err := s.Save(snap); err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestStore_List(t *testing.T) {
	s := newTempStore(t)
	for _, name := range []string{"dev", "staging", "prod"} {
		if err := s.Save(makeSnap(name)); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(names))
	}
}

func TestStore_Delete(t *testing.T) {
	s := newTempStore(t)
	if err := s.Save(makeSnap("dev")); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := s.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Load("dev")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestStore_DeleteMissing(t *testing.T) {
	s := newTempStore(t)
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error deleting nonexistent snapshot, got nil")
	}
}

func TestStore_DefaultDir(t *testing.T) {
	// Verify New with empty dir doesn't error (uses cwd).
	// Clean up the created directory afterwards.
	s, err := store.New("")
	if err != nil {
		t.Fatalf("store.New with empty dir: %v", err)
	}
	_ = s
	_ = os.RemoveAll(".envlock")
}
