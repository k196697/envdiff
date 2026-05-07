package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/patcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestApply_NoMissingKeys(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	results := []diff.Result{
		{Key: "FOO", Status: diff.Match, ValueA: "bar", ValueB: "bar"},
	}
	res, err := patcher.Apply(path, results, patcher.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.KeysAdded) != 0 {
		t.Errorf("expected no keys added, got %v", res.KeysAdded)
	}
}

func TestApply_AppendsMissingKeys(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	results := []diff.Result{
		{Key: "BAR", Status: diff.Missing, ValueA: "baz"},
		{Key: "QUX", Status: diff.Missing, ValueA: ""},
	}
	res, err := patcher.Apply(path, results, patcher.Options{Placeholder: "CHANGEME"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.KeysAdded) != 2 {
		t.Errorf("expected 2 keys added, got %d", len(res.KeysAdded))
	}

	data, _ := os.ReadFile(path)
	body := string(data)
	if !strings.Contains(body, "BAR=baz") {
		t.Errorf("expected BAR=baz in file, got:\n%s", body)
	}
	if !strings.Contains(body, "QUX=CHANGEME") {
		t.Errorf("expected QUX=CHANGEME in file, got:\n%s", body)
	}
}

func TestApply_DryRun_DoesNotModifyFile(t *testing.T) {
	original := "FOO=bar\n"
	path := writeTempEnv(t, original)
	results := []diff.Result{
		{Key: "NEW", Status: diff.Missing, ValueA: "val"},
	}
	res, err := patcher.Apply(path, results, patcher.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if len(res.KeysAdded) != 1 || res.KeysAdded[0] != "NEW" {
		t.Errorf("unexpected KeysAdded: %v", res.KeysAdded)
	}
	data, _ := os.ReadFile(path)
	if string(data) != original {
		t.Errorf("file was modified during dry run")
	}
}

func TestApply_SortedOutput(t *testing.T) {
	path := writeTempEnv(t, "")
	results := []diff.Result{
		{Key: "ZZZ", Status: diff.Missing, ValueA: "1"},
		{Key: "AAA", Status: diff.Missing, ValueA: "2"},
	}
	_, err := patcher.Apply(path, results, patcher.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	body := string(data)
	if strings.Index(body, "AAA") > strings.Index(body, "ZZZ") {
		t.Error("expected keys to be written in sorted order")
	}
}
