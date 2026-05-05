package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoadAll_Basic(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempEnv(t, dir, ".env.prod", "KEY1=val1\nKEY2=val2\n")
	p2 := writeTempEnv(t, dir, ".env.staging", "KEY1=val1\nKEY3=val3\n")

	files, err := loader.LoadAll([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Env["KEY1"] != "val1" {
		t.Errorf("expected KEY1=val1, got %q", files[0].Env["KEY1"])
	}
	if files[1].Env["KEY3"] != "val3" {
		t.Errorf("expected KEY3=val3, got %q", files[1].Env["KEY3"])
	}
}

func TestLoadAll_MissingFile(t *testing.T) {
	_, err := loader.LoadAll([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadAll_PreservesName(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env.test", "FOO=bar\n")
	files, err := loader.LoadAll([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if files[0].Name != ".env.test" {
		t.Errorf("expected name .env.test, got %q", files[0].Name)
	}
}

func TestLoadDir_Basic(t *testing.T) {
	dir := t.TempDir()
	writeTempEnv(t, dir, ".env.prod", "A=1\n")
	writeTempEnv(t, dir, ".env.dev", "A=2\n")

	files, err := loader.LoadDir(dir, ".env.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestLoadDir_NoMatches(t *testing.T) {
	dir := t.TempDir()
	_, err := loader.LoadDir(dir, ".env.*")
	if err == nil {
		t.Error("expected error when no files match, got nil")
	}
}
