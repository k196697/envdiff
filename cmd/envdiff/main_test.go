package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_AllMatch(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	cmd := exec.Command(bin, "--no-color", a, b)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_Mismatch_ExitOne(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=different\n")

	cmd := exec.Command(bin, "--no-color", a, b)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got 0\noutput: %s", out)
	}
}

func TestMain_QuietFlag(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=other\n")

	cmd := exec.Command(bin, "--quiet", a, b)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for mismatch")
	}
	if len(out) != 0 {
		t.Errorf("expected no output with --quiet, got: %s", out)
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit when no args provided")
	}
}
