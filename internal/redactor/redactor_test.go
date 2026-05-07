package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/redactor"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", ValueA: "myapp", ValueB: "myapp", Status: diff.Match},
		{Key: "DB_PASSWORD", ValueA: "hunter2", ValueB: "s3cr3t", Status: diff.Mismatch},
		{Key: "API_KEY", ValueA: "abc123", ValueB: "", Status: diff.MissingInB},
		{Key: "AUTH_TOKEN", ValueA: "", ValueB: "xyz789", Status: diff.MissingInA},
		{Key: "PORT", ValueA: "8080", ValueB: "9090", Status: diff.Mismatch},
	}
}

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	results := makeResults()
	out := redactor.Apply(results, redactor.Options{})

	for _, r := range out {
		switch r.Key {
		case "DB_PASSWORD", "API_KEY", "AUTH_TOKEN":
			if r.ValueA != "***REDACTED***" {
				t.Errorf("key %s: ValueA expected redacted, got %q", r.Key, r.ValueA)
			}
			if r.ValueB != "***REDACTED***" {
				t.Errorf("key %s: ValueB expected redacted, got %q", r.Key, r.ValueB)
			}
		}
	}
}

func TestApply_PreservesNonSensitiveKeys(t *testing.T) {
	results := makeResults()
	out := redactor.Apply(results, redactor.Options{})

	for _, r := range out {
		if r.Key == "APP_NAME" && r.ValueA != "myapp" {
			t.Errorf("APP_NAME should not be redacted, got %q", r.ValueA)
		}
		if r.Key == "PORT" && r.ValueA != "8080" {
			t.Errorf("PORT should not be redacted, got %q", r.ValueA)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults()
	redactor.Apply(results, redactor.Options{})

	for _, r := range results {
		if r.Key == "DB_PASSWORD" && r.ValueA == "***REDACTED***" {
			t.Error("original results were mutated")
		}
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	results := []diff.Result{
		{Key: "STRIPE_KEY", ValueA: "sk_live_abc", ValueB: "sk_live_xyz", Status: diff.Mismatch},
		{Key: "APP_NAME", ValueA: "myapp", ValueB: "myapp", Status: diff.Match},
	}

	out := redactor.Apply(results, redactor.Options{Patterns: []string{"stripe"}})

	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("STRIPE_KEY should be redacted with custom pattern")
	}
	if out[1].ValueA != "myapp" {
		t.Errorf("APP_NAME should not be redacted")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "Db_Password", ValueA: "pass1", ValueB: "pass2", Status: diff.Mismatch},
	}
	out := redactor.Apply(results, redactor.Options{})
	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("expected case-insensitive redaction for Db_Password")
	}
}
