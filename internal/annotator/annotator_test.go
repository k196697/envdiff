package annotator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func writeTempAnnotations(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.annotations")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write annotations file: %v", err)
	}
	return p
}

func TestLoadFile_Basic(t *testing.T) {
	p := writeTempAnnotations(t, "DB_HOST=Database hostname\nDB_PORT=Database port\n")
	anns, err := annotator.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if anns["DB_HOST"] != "Database hostname" {
		t.Errorf("DB_HOST: got %q", anns["DB_HOST"])
	}
	if anns["DB_PORT"] != "Database port" {
		t.Errorf("DB_PORT: got %q", anns["DB_PORT"])
	}
}

func TestLoadFile_CommentsAndBlanks(t *testing.T) {
	p := writeTempAnnotations(t, "# this is a comment\n\nAPI_KEY=Secret API key\n")
	anns, err := annotator.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(anns) != 1 {
		t.Errorf("expected 1 annotation, got %d", len(anns))
	}
	if anns["API_KEY"] != "Secret API key" {
		t.Errorf("API_KEY: got %q", anns["API_KEY"])
	}
}

func TestLoadFile_NonExistent(t *testing.T) {
	anns, err := annotator.LoadFile("/no/such/file.annotations")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(anns) != 0 {
		t.Errorf("expected empty map, got %d entries", len(anns))
	}
}

func TestApply_WithAnnotations(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Match},
		{Key: "API_KEY", Status: diff.Missing},
		{Key: "UNKNOWN", Status: diff.Mismatch},
	}
	anns := map[string]string{
		"DB_HOST": "Database hostname",
		"API_KEY": "Secret API key",
	}
	annotated := annotator.Apply(results, anns)
	if len(annotated) != 3 {
		t.Fatalf("expected 3 results, got %d", len(annotated))
	}
	if annotated[0].Description != "Database hostname" {
		t.Errorf("DB_HOST description: got %q", annotated[0].Description)
	}
	if annotated[1].Description != "Secret API key" {
		t.Errorf("API_KEY description: got %q", annotated[1].Description)
	}
	if annotated[2].Description != "" {
		t.Errorf("UNKNOWN should have empty description, got %q", annotated[2].Description)
	}
}

func TestApply_PreservesResultFields(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Status: diff.Match, ValueA: "8080", ValueB: "8080"},
	}
	annotated := annotator.Apply(results, map[string]string{})
	if annotated[0].Key != "PORT" || annotated[0].ValueA != "8080" {
		t.Errorf("result fields not preserved: %+v", annotated[0])
	}
}
