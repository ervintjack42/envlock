package diff_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/envlock/internal/diff"
	"github.com/user/envlock/internal/snapshot"
)

func makeSnapshot(name string, vars map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Name:      name,
		Vars:      vars,
		CreatedAt: time.Now(),
	}
}

func TestCompare_AddedKey(t *testing.T) {
	from := makeSnapshot("dev", map[string]string{"FOO": "bar"})
	to := makeSnapshot("staging", map[string]string{"FOO": "bar", "NEW_KEY": "value"})

	result := diff.Compare(from, to)

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Type != diff.Added {
		t.Errorf("expected Added, got %s", result.Entries[0].Type)
	}
	if result.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", result.Entries[0].Key)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	from := makeSnapshot("dev", map[string]string{"FOO": "bar", "OLD_KEY": "gone"})
	to := makeSnapshot("staging", map[string]string{"FOO": "bar"})

	result := diff.Compare(from, to)

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Type != diff.Removed {
		t.Errorf("expected Removed, got %s", result.Entries[0].Type)
	}
}

func TestCompare_ChangedKey(t *testing.T) {
	from := makeSnapshot("dev", map[string]string{"DB_URL": "localhost"})
	to := makeSnapshot("prod", map[string]string{"DB_URL": "prod.db.example.com"})

	result := diff.Compare(from, to)

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Type != diff.Changed {
		t.Errorf("expected Changed, got %s", result.Entries[0].Type)
	}
	if result.Entries[0].OldValue != "localhost" {
		t.Errorf("unexpected OldValue: %s", result.Entries[0].OldValue)
	}
}

func TestCompare_NoDifferences(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	from := makeSnapshot("a", vars)
	to := makeSnapshot("b", vars)

	result := diff.Compare(from, to)

	if len(result.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(result.Entries))
	}
}

func TestResult_Format_NoDiff(t *testing.T) {
	from := makeSnapshot("dev", map[string]string{"X": "1"})
	to := makeSnapshot("prod", map[string]string{"X": "1"})

	result := diff.Compare(from, to)
	formatted := result.Format()

	if !strings.Contains(formatted, "No differences") {
		t.Errorf("expected 'No differences' message, got: %s", formatted)
	}
}

func TestResult_Format_WithChanges(t *testing.T) {
	from := makeSnapshot("dev", map[string]string{"PORT": "3000"})
	to := makeSnapshot("prod", map[string]string{"PORT": "8080", "DEBUG": "false"})

	result := diff.Compare(from, to)
	formatted := result.Format()

	if !strings.Contains(formatted, "+") {
		t.Errorf("expected added marker in output: %s", formatted)
	}
	if !strings.Contains(formatted, "~") {
		t.Errorf("expected changed marker in output: %s", formatted)
	}
}
