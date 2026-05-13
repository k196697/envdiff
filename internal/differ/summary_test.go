package differ_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/differ"
)

func makeFileResult(baseline, target string, statuses []diff.Status) differ.FileResult {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: "K", Status: s}
	}
	return differ.FileResult{Baseline: baseline, Target: target, Results: results}
}

func TestSummarise_Clean(t *testing.T) {
	frs := []differ.FileResult{
		makeFileResult("base", "prod", []diff.Status{diff.StatusMatch, diff.StatusMatch}),
	}
	summaries := differ.Summarise(frs)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	s := summaries[0]
	if !s.Clean {
		t.Error("expected Clean=true")
	}
	if s.Match != 2 {
		t.Errorf("expected Match=2, got %d", s.Match)
	}
}

func TestSummarise_WithIssues(t *testing.T) {
	frs := []differ.FileResult{
		makeFileResult("base", "prod", []diff.Status{
			diff.StatusMatch,
			diff.StatusMismatch,
			diff.StatusMissingB,
		}),
	}
	summaries := differ.Summarise(frs)
	s := summaries[0]
	if s.Clean {
		t.Error("expected Clean=false")
	}
	if s.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", s.Mismatch)
	}
	if s.MissingB != 1 {
		t.Errorf("expected MissingB=1, got %d", s.MissingB)
	}
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
}

func TestSummarise_Empty(t *testing.T) {
	summaries := differ.Summarise(nil)
	if len(summaries) != 0 {
		t.Errorf("expected empty summaries, got %d", len(summaries))
	}
}

func TestSummarise_MultiPair(t *testing.T) {
	frs := []differ.FileResult{
		makeFileResult("base", "staging", []diff.Status{diff.StatusMatch}),
		makeFileResult("base", "prod", []diff.Status{diff.StatusMissingA}),
	}
	summaries := differ.Summarise(frs)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	if summaries[0].Target != "staging" {
		t.Errorf("expected staging, got %s", summaries[0].Target)
	}
	if summaries[1].MissingA != 1 {
		t.Errorf("expected MissingA=1, got %d", summaries[1].MissingA)
	}
}
