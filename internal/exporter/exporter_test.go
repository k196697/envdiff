package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/stats"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", File: ".env.prod", Status: diff.StatusMatch},
	{Key: "API_KEY", File: ".env.prod", Status: diff.StatusMismatch},
	{Key: "SECRET", File: ".env.prod", Status: diff.StatusMissing},
}

var sampleStats = stats.Stats{Total: 3, Match: 1, Mismatch: 1, Missing: 1}

func TestExport_TextFormat(t *testing.T) {
	out := filepath.Join(t.TempDir(), "report.txt")
	err := exporter.Export(sampleResults, sampleStats, exporter.Options{
		OutputPath: out,
		Format:     exporter.FormatText,
		Stats:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(out)
	s := string(content)
	if !strings.Contains(s, "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
	if !strings.Contains(s, "Total: 3") {
		t.Error("expected stats in text output")
	}
}

func TestExport_JSONFormat(t *testing.T) {
	out := filepath.Join(t.TempDir(), "report.json")
	err := exporter.Export(sampleResults, sampleStats, exporter.Options{
		OutputPath: out,
		Format:     exporter.FormatJSON,
		Stats:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := payload["results"]; !ok {
		t.Error("expected 'results' key in JSON output")
	}
	if _, ok := payload["stats"]; !ok {
		t.Error("expected 'stats' key in JSON output")
	}
}

func TestExport_MarkdownFormat(t *testing.T) {
	out := filepath.Join(t.TempDir(), "report.md")
	err := exporter.Export(sampleResults, sampleStats, exporter.Options{
		OutputPath: out,
		Format:     exporter.FormatMarkdown,
		Stats:      false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(out)
	s := string(content)
	if !strings.Contains(s, "| Key | File | Status |") {
		t.Error("expected markdown table header")
	}
	if strings.Contains(s, "Total:") {
		t.Error("stats should be omitted when Stats=false")
	}
}

func TestExport_EmptyOutputPath(t *testing.T) {
	err := exporter.Export(sampleResults, sampleStats, exporter.Options{})
	if err == nil {
		t.Fatal("expected error for empty output path")
	}
}

func TestExport_CreatesParentDirs(t *testing.T) {
	out := filepath.Join(t.TempDir(), "nested", "deep", "report.txt")
	err := exporter.Export(sampleResults, sampleStats, exporter.Options{
		OutputPath: out,
		Format:     exporter.FormatText,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}
