package resolver_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/resolver"
)

func TestWriteReport_TextFormat_NoInteresting(t *testing.T) {
	results := []resolver.Result{
		{Key: "FOO", Original: "bar", Resolved: "bar"},
	}
	var buf bytes.Buffer
	if err := resolver.WriteReport(&buf, results, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for unchanged key, got: %q", buf.String())
	}
}

func TestWriteReport_TextFormat_ShowsResolved(t *testing.T) {
	results := []resolver.Result{
		{Key: "URL", Original: "http://${HOST}", Resolved: "http://localhost", Refs: []string{"HOST"}},
	}
	var buf bytes.Buffer
	if err := resolver.WriteReport(&buf, results, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") || !strings.Contains(out, "URL") {
		t.Errorf("expected [OK] line for URL, got: %q", out)
	}
}

func TestWriteReport_TextFormat_ShowsWarning(t *testing.T) {
	results := []resolver.Result{
		{Key: "URL", Original: "http://${HOST}", Resolved: "http://<MISSING:HOST>", Missing: []string{"HOST"}},
	}
	var buf bytes.Buffer
	if err := resolver.WriteReport(&buf, results, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] line, got: %q", out)
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	results := []resolver.Result{
		{Key: "A", Original: "${B}", Resolved: "val", Refs: []string{"B"}},
	}
	var buf bytes.Buffer
	if err := resolver.WriteReport(&buf, results, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []resolver.Result
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 || out[0].Key != "A" {
		t.Errorf("unexpected JSON output: %+v", out)
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	results := []resolver.Result{
		{Key: "X", Original: "${Y}", Resolved: "<MISSING:Y>", Missing: []string{"Y"}},
	}
	var buf bytes.Buffer
	if err := resolver.WriteReport(&buf, results, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Errorf("expected text output, got: %q", buf.String())
	}
}
