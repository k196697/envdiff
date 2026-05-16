package masker_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/masker"
)

var maskedResults = []diff.Result{
	{Key: "DB_PASSWORD", ValueA: masker.DefaultMask, ValueB: masker.DefaultMask, Status: diff.Mismatch},
	{Key: "APP_NAME", ValueA: "myapp", ValueB: "myapp", Status: diff.Match},
	{Key: "PORT", ValueA: "8080", ValueB: "", Status: diff.MissingInB},
}

func TestWriteReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := masker.WriteReport(&buf, maskedResults, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, masker.DefaultMask) {
		t.Error("expected mask string in output")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in output")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := masker.WriteReport(&buf, maskedResults, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rows []struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		ValueA string `json:"value_a"`
		ValueB string `json:"value_b"`
	}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0].ValueA != masker.DefaultMask {
		t.Errorf("expected mask, got %q", rows[0].ValueA)
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	if err := masker.WriteReport(&buf, maskedResults, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Text output should not start with '['
	if strings.HasPrefix(strings.TrimSpace(buf.String()), "[") {
		t.Error("default format should be text, not JSON")
	}
}

func TestWriteReport_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	if err := masker.WriteReport(&buf, []diff.Result{}, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no results, got %q", buf.String())
	}
}
