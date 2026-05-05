package diff

import (
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "8080"}
	b := map[string]string{"HOST": "localhost", "PORT": "8080"}

	results := Compare(a, b)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != StatusMatch {
			t.Errorf("key %q: expected Match, got %v", r.Key, r.Status)
		}
	}
}

func TestCompare_Mismatch(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}

	results := Compare(a, b)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusMismatch {
		t.Errorf("expected Mismatch, got %v", results[0].Status)
	}
	if results[0].ValueA != "8080" || results[0].ValueB != "9090" {
		t.Errorf("unexpected values: A=%q B=%q", results[0].ValueA, results[0].ValueB)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := map[string]string{"SECRET": "abc", "HOST": "localhost"}
	b := map[string]string{"HOST": "localhost"}

	results := Compare(a, b)
	var found bool
	for _, r := range results {
		if r.Key == "SECRET" {
			found = true
			if r.Status != StatusMissingB {
				t.Errorf("expected MissingB for SECRET, got %v", r.Status)
			}
		}
	}
	if !found {
		t.Error("expected result for key SECRET")
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := map[string]string{"HOST": "localhost"}
	b := map[string]string{"HOST": "localhost", "NEW_KEY": "value"}

	results := Compare(a, b)
	var found bool
	for _, r := range results {
		if r.Key == "NEW_KEY" {
			found = true
			if r.Status != StatusMissingA {
				t.Errorf("expected MissingA for NEW_KEY, got %v", r.Status)
			}
		}
	}
	if !found {
		t.Error("expected result for key NEW_KEY")
	}
}

func TestSummary(t *testing.T) {
	results := []Result{
		{Status: StatusMatch},
		{Status: StatusMatch},
		{Status: StatusMismatch},
		{Status: StatusMissingB},
		{Status: StatusMissingA},
		{Status: StatusMissingA},
	}
	match, mismatch, missingB, missingA := Summary(results)
	if match != 2 || mismatch != 1 || missingB != 1 || missingA != 2 {
		t.Errorf("unexpected summary: match=%d mismatch=%d missingB=%d missingA=%d",
			match, mismatch, missingB, missingA)
	}
}
