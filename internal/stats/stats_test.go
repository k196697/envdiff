package stats_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/stats"
)

func makeResults(entries []struct {
	key, val string
	status diff.Status
}) []diff.Result {
	out := make([]diff.Result, len(entries))
	for i, e := range entries {
		out[i] = diff.Result{Key: e.key, Status: e.status}
	}
	return out
}

func TestCompute_AllMatch(t *testing.T) {
	results := map[string][]diff.Result{
		".env.production": {
			{Key: "HOST", Status: diff.StatusMatch},
			{Key: "PORT", Status: diff.StatusMatch},
		},
	}

	s := stats.Compute(results)

	if s.Total != 2 {
		t.Errorf("Total: want 2, got %d", s.Total)
	}
	if s.Matching != 2 {
		t.Errorf("Matching: want 2, got %d", s.Matching)
	}
	if s.Missing != 0 || s.Mismatch != 0 {
		t.Errorf("expected no missing/mismatch")
	}
	if s.Coverage != 100.0 {
		t.Errorf("Coverage: want 100, got %.2f", s.Coverage)
	}
	if !s.IsClean() {
		t.Error("expected IsClean() == true")
	}
}

func TestCompute_WithMissingAndMismatch(t *testing.T) {
	results := map[string][]diff.Result{
		".env.staging": {
			{Key: "HOST", Status: diff.StatusMatch},
			{Key: "SECRET", Status: diff.StatusMissing},
			{Key: "PORT", Status: diff.StatusMismatch},
		},
	}

	s := stats.Compute(results)

	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Missing != 1 {
		t.Errorf("Missing: want 1, got %d", s.Missing)
	}
	if s.Mismatch != 1 {
		t.Errorf("Mismatch: want 1, got %d", s.Mismatch)
	}
	if s.IsClean() {
		t.Error("expected IsClean() == false")
	}
}

func TestCompute_Empty(t *testing.T) {
	s := stats.Compute(map[string][]diff.Result{})

	if s.Total != 0 {
		t.Errorf("Total: want 0, got %d", s.Total)
	}
	if s.Coverage != 0.0 {
		t.Errorf("Coverage: want 0, got %.2f", s.Coverage)
	}
	if !s.IsClean() {
		t.Error("empty result should be clean")
	}
}
