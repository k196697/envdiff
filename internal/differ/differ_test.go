package differ_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/diff"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestRunAll_Basic(t *testing.T) {
	envs := map[string]map[string]string{
		".env.base": {"A": "1", "B": "2"},
		".env.prod": {"A": "1", "B": "3"},
	}
	results, err := differ.RunAll(envs, ".env.base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 FileResult, got %d", len(results))
	}
	fr := results[0]
	if fr.Target != ".env.prod" {
		t.Errorf("expected target .env.prod, got %s", fr.Target)
	}
	var mismatch int
	for _, r := range fr.Results {
		if r.Status == diff.StatusMismatch {
			mismatch++
		}
	}
	if mismatch != 1 {
		t.Errorf("expected 1 mismatch, got %d", mismatch)
	}
}

func TestRunAll_MissingBaseline(t *testing.T) {
	envs := map[string]map[string]string{
		".env.a": {"X": "1"},
	}
	_, err := differ.RunAll(envs, ".env.missing")
	if err == nil {
		t.Fatal("expected error for missing baseline, got nil")
	}
}

func TestRunFiles_Basic(t *testing.T) {
	dir := t.TempDir()
	base := writeTempEnv(t, dir, ".env.base", "A=hello\nB=world\n")
	prod := writeTempEnv(t, dir, ".env.prod", "A=hello\nB=changed\n")

	results, err := differ.RunFiles([]string{base, prod})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Baseline != base {
		t.Errorf("wrong baseline: %s", results[0].Baseline)
	}
}

func TestRunFiles_TooFewFiles(t *testing.T) {
	_, err := differ.RunFiles([]string{"only-one.env"})
	if err == nil {
		t.Fatal("expected error for single file, got nil")
	}
}

func TestRunFiles_MissingFile(t *testing.T) {
	dir := t.TempDir()
	base := writeTempEnv(t, dir, ".env.base", "A=1\n")
	_, err := differ.RunFiles([]string{base, "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing target file, got nil")
	}
}
