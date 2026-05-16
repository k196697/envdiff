package duplicator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/duplicator"
)

func TestDetect_NoDuplicates(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"HOST": "localhost", "PORT": "8080", "DB": "postgres"},
	}
	got := duplicator.Detect(envs)
	if len(got) != 0 {
		t.Fatalf("expected no entries, got %d", len(got))
	}
}

func TestDetect_SingleFileDuplicates(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {
			"DB_PASS":    "secret",
			"REDIS_PASS": "secret",
			"API_KEY":    "secret",
			"HOST":       "localhost",
		},
	}
	got := duplicator.Detect(envs)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	e := got[0]
	if e.Value != "secret" {
		t.Errorf("expected value 'secret', got %q", e.Value)
	}
	if len(e.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(e.Keys), e.Keys)
	}
	if e.File != "prod" {
		t.Errorf("expected file 'prod', got %q", e.File)
	}
}

func TestDetect_EmptyValuesIgnored(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"A": "", "B": "", "C": "real"},
	}
	got := duplicator.Detect(envs)
	if len(got) != 0 {
		t.Fatalf("empty values should not be flagged, got %d entries", len(got))
	}
}

func TestDetect_MultipleFiles(t *testing.T) {
	envs := map[string]map[string]string{
		"alpha": {"X": "42", "Y": "42"},
		"beta":  {"P": "hello", "Q": "world"},
	}
	got := duplicator.Detect(envs)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].File != "alpha" {
		t.Errorf("expected file 'alpha', got %q", got[0].File)
	}
}

func TestDetect_SortedOutput(t *testing.T) {
	envs := map[string]map[string]string{
		"z": {"K1": "dup", "K2": "dup"},
		"a": {"M1": "dup", "M2": "dup"},
	}
	got := duplicator.Detect(envs)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	// Files should be sorted: "a" before "z".
	if got[0].File != "a" || got[1].File != "z" {
		t.Errorf("files not sorted: %q %q", got[0].File, got[1].File)
	}
}

func TestWrite_Clean(t *testing.T) {
	var buf bytes.Buffer
	duplicator.Write(&buf, nil)
	if !strings.Contains(buf.String(), "no duplicate") {
		t.Errorf("expected clean message, got %q", buf.String())
	}
}

func TestWrite_ShowsEntries(t *testing.T) {
	entries := []duplicator.Entry{
		{File: "prod", Value: "secret", Keys: []string{"DB_PASS", "REDIS_PASS"}},
	}
	var buf bytes.Buffer
	duplicator.Write(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("output missing file name: %q", out)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("output missing value: %q", out)
	}
	if !strings.Contains(out, "DB_PASS") || !strings.Contains(out, "REDIS_PASS") {
		t.Errorf("output missing key names: %q", out)
	}
}
