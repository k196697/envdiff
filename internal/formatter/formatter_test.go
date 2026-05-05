package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch, ValueA: "myapp", ValueB: "myapp"},
		{Key: "DB_HOST", Status: diff.StatusMismatch, ValueA: "localhost", ValueB: "prod-db"},
		{Key: "SECRET", Status: diff.StatusMissingInB, ValueA: "abc123", ValueB: ""},
		{Key: "NEW_KEY", Status: diff.StatusMissingInA, ValueA: "", ValueB: "value"},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.Write(&buf, sampleResults(), formatter.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[=] APP_NAME") {
		t.Errorf("expected match line, got:\n%s", out)
	}
	if !strings.Contains(out, "[~] DB_HOST") {
		t.Errorf("expected mismatch line, got:\n%s", out)
	}
	if !strings.Contains(out, "[-] SECRET") {
		t.Errorf("expected missing-in-B line, got:\n%s", out)
	}
	if !strings.Contains(out, "[+] NEW_KEY") {
		t.Errorf("expected missing-in-A line, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.Write(&buf, sampleResults(), formatter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\noutput: %s", err, buf.String())
	}
	if result["match_count"].(float64) != 1 {
		t.Errorf("expected match_count=1, got %v", result["match_count"])
	}
	mismatched, ok := result["mismatched"].(map[string]interface{})
	if !ok {
		t.Fatal("expected mismatched map in JSON output")
	}
	if _, ok := mismatched["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST in mismatched")
	}
}

func TestWrite_JSONFormat_AllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Status: diff.StatusMatch, ValueA: "bar", ValueB: "bar"},
	}
	var buf bytes.Buffer
	err := formatter.Write(&buf, results, formatter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, exists := result["mismatched"]; exists {
		t.Errorf("mismatched should be omitted when empty")
	}
}

func TestWrite_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.Write(&buf, sampleResults(), formatter.Format("unknown"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[=") {
		t.Errorf("expected text format fallback")
	}
}
