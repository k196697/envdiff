package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestReport_AllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", ValueA: "8080", ValueB: "8080", Status: diff.StatusMatch},
	}
	var buf bytes.Buffer
	Report(results, "dev.env", "prod.env", &buf)
	out := buf.String()

	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1 match") {
		t.Errorf("expected '1 match' in summary, got:\n%s", out)
	}
}

func TestReport_Mismatch(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", ValueA: "localhost", ValueB: "db.prod.example.com", Status: diff.StatusMismatch},
	}
	var buf bytes.Buffer
	Report(results, "dev.env", "prod.env", &buf)
	out := buf.String()

	if !strings.Contains(out, "[MISMATCH]") {
		t.Errorf("expected [MISMATCH] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected value 'localhost' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1 mismatch") {
		t.Errorf("expected '1 mismatch' in summary, got:\n%s", out)
	}
}

func TestReport_Missing(t *testing.T) {
	results := []diff.Result{
		{Key: "SECRET_KEY", ValueA: "abc123", ValueB: "", Status: diff.StatusMissingInB},
		{Key: "NEW_FLAG", ValueA: "", ValueB: "true", Status: diff.StatusMissingInA},
	}
	var buf bytes.Buffer
	Report(results, "dev.env", "prod.env", &buf)
	out := buf.String()

	if !strings.Contains(out, "[MISSING]") {
		t.Errorf("expected [MISSING] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2 missing") {
		t.Errorf("expected '2 missing' in summary, got:\n%s", out)
	}
}

func TestReport_SortedOutput(t *testing.T) {
	results := []diff.Result{
		{Key: "ZEBRA", ValueA: "z", ValueB: "z", Status: diff.StatusMatch},
		{Key: "ALPHA", ValueA: "a", ValueB: "a", Status: diff.StatusMatch},
		{Key: "MANGO", ValueA: "m", ValueB: "m", Status: diff.StatusMatch},
	}
	var buf bytes.Buffer
	Report(results, "a.env", "b.env", &buf)
	out := buf.String()

	alpha := strings.Index(out, "ALPHA")
	mango := strings.Index(out, "MANGO")
	zebra := strings.Index(out, "ZEBRA")

	if !(alpha < mango && mango < zebra) {
		t.Errorf("expected keys in sorted order, got:\n%s", out)
	}
}
