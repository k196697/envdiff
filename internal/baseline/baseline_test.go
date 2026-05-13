package baseline_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/envdiff/internal/baseline"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := baseline.Save(path, ".env.production", env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if snap.Source != ".env.production" {
		t.Errorf("source: got %q, want %q", snap.Source, ".env.production")
	}
	if snap.Keys["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", snap.Keys["APP_ENV"], "production")
	}
	if snap.Keys["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", snap.Keys["PORT"], "8080")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "snap.json")

	if err := baseline.Save(path, "src", map[string]string{"K": "v"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestDiff_Added(t *testing.T) {
	snap := &baseline.Snapshot{Keys: map[string]string{"A": "1"}}
	live := map[string]string{"A": "1", "B": "2"}
	added, removed, changed := baseline.Diff(snap, live)
	if len(added) != 1 || added[0] != "B" {
		t.Errorf("added: got %v, want [B]", added)
	}
	if len(removed) != 0 {
		t.Errorf("removed: got %v, want []", removed)
	}
	if len(changed) != 0 {
		t.Errorf("changed: got %v, want []", changed)
	}
}

func TestDiff_RemovedAndChanged(t *testing.T) {
	snap := &baseline.Snapshot{Keys: map[string]string{"A": "1", "B": "old", "C": "3"}}
	live := map[string]string{"A": "1", "B": "new"}
	added, removed, changed := baseline.Diff(snap, live)

	if len(added) != 0 {
		t.Errorf("added: got %v, want []", added)
	}
	sort.Strings(removed)
	if len(removed) != 1 || removed[0] != "C" {
		t.Errorf("removed: got %v, want [C]", removed)
	}
	if len(changed) != 1 || changed[0] != "B" {
		t.Errorf("changed: got %v, want [B]", changed)
	}
}

func TestDiff_AllMatch(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	snap := &baseline.Snapshot{Keys: env}
	added, removed, changed := baseline.Diff(snap, map[string]string{"X": "1", "Y": "2"})
	if len(added)+len(removed)+len(changed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v", added, removed, changed)
	}
}
