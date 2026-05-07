package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/ignorer"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempIgnore: %v", err)
	}
	return path
}

func TestLoad_NonExistentFile(t *testing.T) {
	ig, err := ignorer.Load("/nonexistent/.envdiffignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if ig.Match("ANYTHING") {
		t.Error("empty ignorer should not match any key")
	}
}

func TestLoad_CommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "# this is a comment\n\nSECRET_KEY\n")
	ig, err := ignorer.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !ig.Match("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if ig.Match("OTHER_KEY") {
		t.Error("OTHER_KEY should not be ignored")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	path := writeTempIgnore(t, "secret_key\n")
	ig, err := ignorer.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, key := range []string{"SECRET_KEY", "secret_key", "Secret_Key"} {
		if !ig.Match(key) {
			t.Errorf("expected %q to match (case-insensitive)", key)
		}
	}
}

func TestFilter_RemovesIgnoredKeys(t *testing.T) {
	path := writeTempIgnore(t, "AWS_SECRET\nDB_PASS\n")
	ig, err := ignorer.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	input := []string{"AWS_SECRET", "APP_ENV", "DB_PASS", "PORT"}
	got := ig.Filter(input)
	want := []string{"APP_ENV", "PORT"}
	if len(got) != len(want) {
		t.Fatalf("Filter: got %v, want %v", got, want)
	}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("Filter[%d]: got %q, want %q", i, got[i], k)
		}
	}
}

func TestFilter_NoPatterns(t *testing.T) {
	ig, _ := ignorer.Load("/nonexistent/.envdiffignore")
	input := []string{"A", "B", "C"}
	got := ig.Filter(input)
	if len(got) != len(input) {
		t.Errorf("expected all keys preserved, got %v", got)
	}
}
