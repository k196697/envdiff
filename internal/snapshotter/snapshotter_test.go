package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshotter"
)

func TestCapture_CopiesEnv(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshotter.Capture("test", env)
	env["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected original value, got %s", s.Env["FOO"])
	}
	if s.Name != "test" {
		t.Errorf("unexpected name: %s", s.Name)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	env := map[string]string{"KEY": "value", "OTHER": "123"}
	orig := snapshotter.Capture("prod", env)

	if err := snapshotter.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := snapshotter.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != orig.Name {
		t.Errorf("name mismatch: %s vs %s", loaded.Name, orig.Name)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("unexpected value: %s", loaded.Env["KEY"])
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "snap.json")
	s := snapshotter.Capture("x", map[string]string{})
	if err := snapshotter.Save(s, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestLoad_NonExistent(t *testing.T) {
	_, err := snapshotter.Load("/no/such/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDiff_Added(t *testing.T) {
	before := snapshotter.Capture("before", map[string]string{"A": "1"})
	after := snapshotter.Capture("after", map[string]string{"A": "1", "B": "2"})
	deltas := snapshotter.Diff(before, after)
	if len(deltas) != 1 || deltas[0].Kind != snapshotter.Added || deltas[0].Key != "B" {
		t.Errorf("unexpected deltas: %+v", deltas)
	}
}

func TestDiff_Removed(t *testing.T) {
	before := snapshotter.Capture("before", map[string]string{"A": "1", "B": "2"})
	after := snapshotter.Capture("after", map[string]string{"A": "1"})
	deltas := snapshotter.Diff(before, after)
	if len(deltas) != 1 || deltas[0].Kind != snapshotter.Removed || deltas[0].Key != "B" {
		t.Errorf("unexpected deltas: %+v", deltas)
	}
}

func TestDiff_Changed(t *testing.T) {
	before := snapshotter.Capture("before", map[string]string{"A": "old"})
	after := snapshotter.Capture("after", map[string]string{"A": "new"})
	deltas := snapshotter.Diff(before, after)
	if len(deltas) != 1 || deltas[0].Kind != snapshotter.Changed {
		t.Errorf("unexpected deltas: %+v", deltas)
	}
	if deltas[0].OldValue != "old" || deltas[0].NewValue != "new" {
		t.Errorf("wrong values: %+v", deltas[0])
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	before := snapshotter.Capture("a", env)
	after := snapshotter.Capture("b", env)
	deltas := snapshotter.Diff(before, after)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas, got %+v", deltas)
	}
}
