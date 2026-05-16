package pinner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/pinner"
)

func writeTempPins(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "pinned.txt")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_EmptyPath(t *testing.T) {
	ps, err := pinner.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.IsPinned("FOO") {
		t.Error("expected no pinned keys")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	ps, err := pinner.Load("/no/such/file.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.IsPinned("ANY") {
		t.Error("expected empty set")
	}
}

func TestLoad_CommentsAndBlanks(t *testing.T) {
	p := writeTempPins(t, "# comment\n\nDB_HOST\nAPI_KEY\n")
	ps, err := pinner.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ps.IsPinned("DB_HOST") {
		t.Error("expected DB_HOST to be pinned")
	}
	if !ps.IsPinned("API_KEY") {
		t.Error("expected API_KEY to be pinned")
	}
	if ps.IsPinned("OTHER") {
		t.Error("OTHER should not be pinned")
	}
}

func TestIsPinned_CaseInsensitive(t *testing.T) {
	p := writeTempPins(t, "db_host\n")
	ps, _ := pinner.Load(p)
	if !ps.IsPinned("DB_HOST") {
		t.Error("expected case-insensitive match")
	}
}

func TestApply_RemovesMismatchForPinnedKey(t *testing.T) {
	p := writeTempPins(t, "DB_HOST\n")
	ps, _ := pinner.Load(p)

	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Mismatch},
		{Key: "API_KEY", Status: diff.Mismatch},
		{Key: "PORT", Status: diff.Match},
	}
	out := ps.Apply(results)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Key == "DB_HOST" {
			t.Error("DB_HOST mismatch should have been filtered")
		}
	}
}

func TestApply_PreservesMissingForPinnedKey(t *testing.T) {
	p := writeTempPins(t, "DB_HOST\n")
	ps, _ := pinner.Load(p)

	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.MissingInB},
	}
	out := ps.Apply(results)
	if len(out) != 1 {
		t.Errorf("expected missing entry to be preserved, got %d results", len(out))
	}
}
