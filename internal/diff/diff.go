package diff

import (
	"fmt"
	"sort"

	"github.com/user/envlock/internal/snapshot"
)

// ChangeType represents the kind of change between two snapshots.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Entry represents a single diff entry between two snapshots.
type Entry struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two snapshots.
type Result struct {
	From    string
	To      string
	Entries []Entry
}

// Compare computes the diff between two snapshots.
func Compare(from, to *snapshot.Snapshot) *Result {
	result := &Result{
		From: from.Name,
		To:   to.Name,
	}

	for key, newVal := range to.Vars {
		if oldVal, exists := from.Vars[key]; !exists {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				Type:     Added,
				NewValue: newVal,
			})
		} else if oldVal != newVal {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				Type:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for key, oldVal := range from.Vars {
		if _, exists := to.Vars[key]; !exists {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				Type:     Removed,
				OldValue: oldVal,
			})
		}
	}

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Key < result.Entries[j].Key
	})

	return result
}

// Format returns a human-readable string representation of the diff.
func (r *Result) Format() string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("No differences between '%s' and '%s'.\n", r.From, r.To)
	}

	out := fmt.Sprintf("Diff: %s → %s\n", r.From, r.To)
	for _, e := range r.Entries {
		switch e.Type {
		case Added:
			out += fmt.Sprintf("  + %s=%s\n", e.Key, e.NewValue)
		case Removed:
			out += fmt.Sprintf("  - %s=%s\n", e.Key, e.OldValue)
		case Changed:
			out += fmt.Sprintf("  ~ %s: %s → %s\n", e.Key, e.OldValue, e.NewValue)
		}
	}
	return out
}
