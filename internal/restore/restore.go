// Package restore applies a saved snapshot back to the current environment
// by writing an export script that can be sourced by the user's shell.
package restore

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envlock/internal/snapshot"
)

// ShellFormat controls the output format of the restore script.
type ShellFormat string

const (
	FormatBash ShellFormat = "bash"
	FormatFish ShellFormat = "fish"
)

// Options configures the restore operation.
type Options struct {
	// Format is the target shell format (default: bash).
	Format ShellFormat
	// Out is the writer to emit the script to (default: os.Stdout).
	Out io.Writer
	// Unset removes variables that are not present in the snapshot.
	Unset bool
}

// WriteScript writes a shell-sourceable export script for the given snapshot.
func WriteScript(snap *snapshot.Snapshot, opts Options) error {
	if snap == nil {
		return fmt.Errorf("restore: snapshot must not be nil")
	}

	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.Format == "" {
		opts.Format = FormatBash
	}

	switch opts.Format {
	case FormatBash:
		return writeBash(snap, opts)
	case FormatFish:
		return writeFish(snap, opts)
	default:
		return fmt.Errorf("restore: unsupported shell format %q", opts.Format)
	}
}

func writeBash(snap *snapshot.Snapshot, opts Options) error {
	fmt.Fprintf(opts.Out, "# envlock restore — snapshot: %s\n", snap.Name)
	for k, v := range snap.Vars {
		escaped := strings.ReplaceAll(v, "'", "'\"'\"'")
		fmt.Fprintf(opts.Out, "export %s='%s'\n", k, escaped)
	}
	if opts.Unset {
		for _, pair := range os.Environ() {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) < 1 {
				continue
			}
			key := parts[0]
			if _, exists := snap.Vars[key]; !exists {
				fmt.Fprintf(opts.Out, "unset %s\n", key)
			}
		}
	}
	return nil
}

func writeFish(snap *snapshot.Snapshot, opts Options) error {
	fmt.Fprintf(opts.Out, "# envlock restore — snapshot: %s\n", snap.Name)
	for k, v := range snap.Vars {
		escaped := strings.ReplaceAll(v, "'", "\'")
		fmt.Fprintf(opts.Out, "set -x %s '%s'\n", k, escaped)
	}
	if opts.Unset {
		for _, pair := range os.Environ() {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) < 1 {
				continue
			}
			key := parts[0]
			if _, exists := snap.Vars[key]; !exists {
				fmt.Fprintf(opts.Out, "set -e %s\n", key)
			}
		}
	}
	return nil
}
