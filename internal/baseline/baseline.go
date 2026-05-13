// Package baseline provides functionality for pinning and comparing
// a reference .env snapshot against current environment files.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a pinned state of environment variables.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Keys      map[string]string `json:"keys"`
}

// Save writes a snapshot of the given env map to a JSON file at path.
func Save(path string, source string, env map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("baseline: create dirs: %w", err)
	}

	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Keys:      env,
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("baseline: encode snapshot: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot from a JSON file at path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open file: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("baseline: decode snapshot: %w", err)
	}
	return &snap, nil
}

// Diff compares a live env map against a snapshot and returns added,
// removed, and changed keys.
func Diff(snap *Snapshot, live map[string]string) (added, removed, changed []string) {
	for k := range live {
		if _, ok := snap.Keys[k]; !ok {
			added = append(added, k)
		}
	}
	for k, sv := range snap.Keys {
		lv, ok := live[k]
		if !ok {
			removed = append(removed, k)
		} else if lv != sv {
			changed = append(changed, k)
		}
	}
	return
}
