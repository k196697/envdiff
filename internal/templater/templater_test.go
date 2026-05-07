package templater_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/templater"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Status: diff.Match, ValueA: "production", ValueB: "production"},
		{Key: "DB_HOST", Status: diff.Mismatch, ValueA: "localhost", ValueB: "db.prod"},
		{Key: "SECRET_KEY", Status: diff.MissingInB, ValueA: "abc123", ValueB: ""},
		{Key: "LOG_LEVEL", Status: diff.MissingInA, ValueA: "", ValueB: "debug"},
	}
}

func TestGenerate_DefaultPlaceholder(t *testing.T) {
	var buf strings.Builder
	err := templater.Generate(&buf, sampleResults(), templater.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=REPLACE_ME") {
		t.Errorf("expected DB_HOST placeholder, got:\n%s", out)
	}
	if !strings.Contains(out, "SECRET_KEY=REPLACE_ME") {
		t.Errorf("expected SECRET_KEY placeholder, got:\n%s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL=REPLACE_ME") {
		t.Errorf("expected LOG_LEVEL placeholder, got:\n%s", out)
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	var buf strings.Builder
	err := templater.Generate(&buf, sampleResults(), templater.Options{Placeholder: "TODO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=TODO") {
		t.Errorf("expected DB_HOST=TODO, got:\n%s", out)
	}
}

func TestGenerate_IncludeValues(t *testing.T) {
	var buf strings.Builder
	err := templater.Generate(&buf, sampleResults(), templater.Options{IncludeValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production, got:\n%s", out)
	}
	// Mismatched key should still use placeholder
	if !strings.Contains(out, "DB_HOST=REPLACE_ME") {
		t.Errorf("expected DB_HOST=REPLACE_ME, got:\n%s", out)
	}
}

func TestGenerate_SortedOutput(t *testing.T) {
	var buf strings.Builder
	_ = templater.Generate(&buf, sampleResults(), templater.Options{})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		prev := strings.Split(lines[i-1], "=")[0]
		curr := strings.Split(lines[i], "=")[0]
		if prev > curr {
			t.Errorf("output not sorted: %q comes before %q", prev, curr)
		}
	}
}

func TestGenerate_Empty(t *testing.T) {
	var buf strings.Builder
	err := templater.Generate(&buf, []diff.Result{}, templater.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}
