package classifier_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/diff"
)

func sampleClassified() []classifier.Result {
	return []classifier.Result{
		{Diff: diff.Result{Key: "SECRET_KEY", Status: "missing"}, Severity: classifier.SeverityCritical},
		{Diff: diff.Result{Key: "DB_HOST", Status: "mismatch"}, Severity: classifier.SeverityWarning},
		{Diff: diff.Result{Key: "APP_NAME", Status: "match"}, Severity: classifier.SeverityInfo},
	}
}

func TestWriteReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := classifier.WriteReport(&buf, sampleClassified(), "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "CRITICAL") {
		t.Error("expected CRITICAL in output")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING in output")
	}
	if !strings.Contains(out, "INFO") {
		t.Error("expected INFO in output")
	}
}

func TestWriteReport_TextFormat_CriticalFirst(t *testing.T) {
	var buf bytes.Buffer
	_ = classifier.WriteReport(&buf, sampleClassified(), "text")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.Contains(lines[0], "CRITICAL") {
		t.Errorf("first line should be CRITICAL, got: %s", lines[0])
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := classifier.WriteReport(&buf, sampleClassified(), "json"); err != nil {
		t.Fatal(err)
	}
	var rows []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if r["severity"] == "" {
			t.Errorf("missing severity field in row %v", r)
		}
	}
}

func TestWriteReport_EmptyResults_Text(t *testing.T) {
	var buf bytes.Buffer
	_ = classifier.WriteReport(&buf, nil, "text")
	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results' message for empty input")
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	_ = classifier.WriteReport(&buf, sampleClassified(), "")
	if !strings.Contains(buf.String(), "CRITICAL") {
		t.Error("default format should be text")
	}
}
