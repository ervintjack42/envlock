package restore_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/envlock/internal/restore"
	"github.com/user/envlock/internal/snapshot"
)

func makeSnap(name string, vars map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Name:      name,
		Vars:      vars,
		CreatedAt: time.Now(),
	}
}

func TestWriteScript_BashFormat(t *testing.T) {
	snap := makeSnap("dev", map[string]string{
		"APP_ENV": "development",
		"PORT":    "8080",
	})

	var buf strings.Builder
	err := restore.WriteScript(snap, restore.Options{
		Format: restore.FormatBash,
		Out:    &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "export APP_ENV='development'") {
		t.Errorf("expected bash export for APP_ENV, got:\n%s", output)
	}
	if !strings.Contains(output, "export PORT='8080'") {
		t.Errorf("expected bash export for PORT, got:\n%s", output)
	}
	if !strings.Contains(output, "# envlock restore — snapshot: dev") {
		t.Errorf("expected header comment, got:\n%s", output)
	}
}

func TestWriteScript_FishFormat(t *testing.T) {
	snap := makeSnap("staging", map[string]string{
		"DB_URL": "postgres://localhost/staging",
	})

	var buf strings.Builder
	err := restore.WriteScript(snap, restore.Options{
		Format: restore.FormatFish,
		Out:    &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "set -x DB_URL") {
		t.Errorf("expected fish set -x for DB_URL, got:\n%s", output)
	}
}

func TestWriteScript_NilSnapshotReturnsError(t *testing.T) {
	err := restore.WriteScript(nil, restore.Options{})
	if err == nil {
		t.Fatal("expected error for nil snapshot, got nil")
	}
}

func TestWriteScript_UnsupportedFormatReturnsError(t *testing.T) {
	snap := makeSnap("dev", map[string]string{"K": "V"})
	err := restore.WriteScript(snap, restore.Options{
		Format: "zsh",
		Out:    &strings.Builder{},
	})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestWriteScript_EscapesSingleQuotes(t *testing.T) {
	snap := makeSnap("prod", map[string]string{
		"MSG": "it's alive",
	})

	var buf strings.Builder
	err := restore.WriteScript(snap, restore.Options{
		Format: restore.FormatBash,
		Out:    &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "'it's alive'") {
		t.Errorf("single quote not escaped properly in output:\n%s", output)
	}
}
