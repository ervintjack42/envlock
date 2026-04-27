package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_ParsesEnviron(t *testing.T) {
	environ := []string{
		"HOME=/root",
		"PATH=/usr/bin:/bin",
		"EMPTY=",
		"WITH_EQUALS=a=b=c",
	}

	s, err := New("test", environ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"HOME":        "/root",
		"PATH":        "/usr/bin:/bin",
		"EMPTY":       "",
		"WITH_EQUALS": "a=b=c",
	}
	for key, want := range cases {
		if got := s.Env[key]; got != want {
			t.Errorf("Env[%q] = %q, want %q", key, got, want)
		}
	}
}

func TestNew_EmptyNameReturnsError(t *testing.T) {
	_, err := New("", []string{"FOO=bar"})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	environ := []string{"FOO=bar", "BAZ=qux"}
	orig, err := New("roundtrip", environ)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Name != orig.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, orig.Name)
	}
	for k, v := range orig.Env {
		if loaded.Env[k] != v {
			t.Errorf("Env[%q]: got %q, want %q", k, loaded.Env[k], v)
		}
	}
}

func TestLoad_MissingFileReturnsError(t *testing.T) {
	_, err := Load(filepath.Join(os.TempDir(), "nonexistent_envlock_snap.json"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
