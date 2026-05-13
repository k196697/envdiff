// Package snapshotter captures point-in-time snapshots of parsed env files
// and allows comparison against a previous snapshot to detect drift.
package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot holds a named capture of env key/value pairs.
type Snapshot struct {
	Name      string            `json:"name"`
	CapturedAt time.Time        `json:"captured_at"`
	Env       map[string]string `json:"env"`
}

// Capture creates a new Snapshot from the provided env map.
func Capture(name string, env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		Name:       name,
		CapturedAt: time.Now().UTC(),
		Env:        copy,
	}
}

// Save writes a snapshot to disk as JSON.
func Save(s Snapshot, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshotter: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshotter: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from disk.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: read: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	return s, nil
}

// DeltaKind describes how a key changed between two snapshots.
type DeltaKind string

const (
	Added   DeltaKind = "added"
	Removed DeltaKind = "removed"
	Changed DeltaKind = "changed"
)

// Delta represents a single key-level difference between two snapshots.
type Delta struct {
	Key      string    `json:"key"`
	Kind     DeltaKind `json:"kind"`
	OldValue string    `json:"old_value,omitempty"`
	NewValue string    `json:"new_value,omitempty"`
}

// Diff compares two snapshots and returns the list of deltas.
func Diff(before, after Snapshot) []Delta {
	var deltas []Delta
	for k, newVal := range after.Env {
		if oldVal, ok := before.Env[k]; !ok {
			deltas = append(deltas, Delta{Key: k, Kind: Added, NewValue: newVal})
		} else if oldVal != newVal {
			deltas = append(deltas, Delta{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}
	for k, oldVal := range before.Env {
		if _, ok := after.Env[k]; !ok {
			deltas = append(deltas, Delta{Key: k, Kind: Removed, OldValue: oldVal})
		}
	}
	return deltas
}
