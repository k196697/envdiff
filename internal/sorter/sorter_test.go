package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", FileA: "prod", FileB: "dev", ValueA: "z", ValueB: "z"},
		{Key: "ALPHA", FileA: "prod", FileB: "dev", Mismatch: true, ValueA: "1", ValueB: "2"},
		{Key: "MISSING_B", FileA: "prod", FileB: "dev", MissingInB: true},
		{Key: "BETA", FileA: "prod", FileB: "dev", ValueA: "b", ValueB: "b"},
		{Key: "MISSING_A", FileA: "prod", FileB: "dev", MissingInA: true},
	}
}

func TestSort_ByKey(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.ByKey)

	expected := []string{"ALPHA", "BETA", "MISSING_A", "MISSING_B", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByStatus(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.ByStatus)

	// Missing entries should come first, then mismatches, then matches.
	if !sorted[0].MissingInA && !sorted[0].MissingInB {
		t.Errorf("expected first result to be missing, got key %q", sorted[0].Key)
	}
	if !sorted[1].MissingInA && !sorted[1].MissingInB {
		t.Errorf("expected second result to be missing, got key %q", sorted[1].Key)
	}
	if !sorted[2].Mismatch {
		t.Errorf("expected third result to be mismatch, got key %q", sorted[2].Key)
	}
}

func TestSort_ByFile(t *testing.T) {
	results := []diff.Result{
		{Key: "Z", FileA: "staging"},
		{Key: "A", FileA: "prod"},
		{Key: "M", FileA: "prod"},
	}
	sorted := sorter.Sort(results, sorter.ByFile)

	if sorted[0].FileA != "prod" {
		t.Errorf("expected first FileA to be prod, got %q", sorted[0].FileA)
	}
	if sorted[0].Key != "A" {
		t.Errorf("expected first key to be A, got %q", sorted[0].Key)
	}
	if sorted[2].FileA != "staging" {
		t.Errorf("expected last FileA to be staging, got %q", sorted[2].FileA)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	results := sampleResults()
	originalFirst := results[0].Key

	sorter.Sort(results, sorter.ByKey)

	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated: first key changed from %q to %q", originalFirst, results[0].Key)
	}
}
