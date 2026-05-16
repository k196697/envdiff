package masker_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/masker"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_PASSWORD", ValueA: "hunter2", ValueB: "s3cr3t", Status: diff.Mismatch},
		{Key: "API_KEY", ValueA: "abc123", ValueB: "", Status: diff.MissingInB},
		{Key: "APP_NAME", ValueA: "myapp", ValueB: "myapp", Status: diff.Match},
		{Key: "PORT", ValueA: "8080", ValueB: "9090", Status: diff.Mismatch},
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	results := makeResults()
	out := masker.Apply(results, masker.Options{})

	if out[0].ValueA != masker.DefaultMask {
		t.Errorf("DB_PASSWORD ValueA: got %q, want %q", out[0].ValueA, masker.DefaultMask)
	}
	if out[0].ValueB != masker.DefaultMask {
		t.Errorf("DB_PASSWORD ValueB: got %q, want %q", out[0].ValueB, masker.DefaultMask)
	}
}

func TestApply_PreservesNonSensitiveValues(t *testing.T) {
	out := masker.Apply(makeResults(), masker.Options{})

	if out[2].ValueA != "myapp" {
		t.Errorf("APP_NAME ValueA: got %q, want %q", out[2].ValueA, "myapp")
	}
	if out[3].ValueB != "9090" {
		t.Errorf("PORT ValueB: got %q, want %q", out[3].ValueB, "9090")
	}
}

func TestApply_DoesNotMaskEmptyValues(t *testing.T) {
	out := masker.Apply(makeResults(), masker.Options{})

	// API_KEY is missing in B, so ValueB is empty — must stay empty.
	if out[1].ValueB != "" {
		t.Errorf("API_KEY ValueB should remain empty, got %q", out[1].ValueB)
	}
}

func TestApply_CustomMaskString(t *testing.T) {
	out := masker.Apply(makeResults(), masker.Options{Mask: "[REDACTED]"})

	if out[0].ValueA != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out[0].ValueA)
	}
}

func TestApply_ExtraPatterns(t *testing.T) {
	results := []diff.Result{
		{Key: "STRIPE_WEBHOOK", ValueA: "whsec_abc", ValueB: "whsec_xyz", Status: diff.Mismatch},
		{Key: "PLAIN_VAR", ValueA: "hello", ValueB: "hello", Status: diff.Match},
	}
	out := masker.Apply(results, masker.Options{ExtraPatterns: []string{"webhook"}})

	if out[0].ValueA != masker.DefaultMask {
		t.Errorf("STRIPE_WEBHOOK should be masked, got %q", out[0].ValueA)
	}
	if out[1].ValueA != "hello" {
		t.Errorf("PLAIN_VAR should not be masked, got %q", out[1].ValueA)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := makeResults()
	masker.Apply(original, masker.Options{})

	if original[0].ValueA != "hunter2" {
		t.Error("Apply mutated the original slice")
	}
}
