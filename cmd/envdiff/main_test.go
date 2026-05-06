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

func runEnvdiff(t *testing.T, bin string, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	return cmd.CombinedOutput()
}

func TestMain_AllMatch(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	out, err := runEnvdiff(t, bin, "--no-color", a, b)
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_Mismatch_ExitOne(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=different\n")

	out, err := runEnvdiff(t, bin, "--no-color", a, b)
	if err == nil {
		t.Fatalf("expected non-zero exit, got 0\noutput: %s", out)
	}
}

func TestMain_QuietFlag(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=other\n")

	out, err := runEnvdiff(t, bin, "--quiet", a, b)
	if err == nil {
		t.Fatalf("expected non-zero exit for mismatch")
	}
	if len(out) != 0 {
		t.Errorf("expected no output with --quiet, got: %s", out)
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	out, err := runEnvdiff(t, bin)
	if err == nil {
		t.Fatalf("expected non-zero exit when no args provided\noutput: %s", out)
	}
}
