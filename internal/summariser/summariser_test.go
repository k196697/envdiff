package summariser_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summariser"
)

func makeFileResults() []summariser.FileResult {
	return []summariser.FileResult{
		{
			Name: "staging.env",
			Results: []diff.Result{
				{Key: "DB_HOST", Status: diff.Match},
				{Key: "API_KEY", Status: diff.Mismatch},
				{Key: "SECRET", Status: diff.MissingInB},
			},
		},
		{
			Name: "production.env",
			Results: []diff.Result{
				{Key: "DB_HOST", Status: diff.Match},
				{Key: "DB_PASS", Status: diff.Match},
			},
		},
	}
}

func TestBuild_Totals(t *testing.T) {
	s := summariser.Build(makeFileResults())

	if s.TotalFiles != 2 {
		t.Errorf("TotalFiles: got %d, want 2", s.TotalFiles)
	}
	if s.TotalKeys != 5 {
		t.Errorf("TotalKeys: got %d, want 5", s.TotalKeys)
	}
	if s.Matched != 3 {
		t.Errorf("Matched: got %d, want 3", s.Matched)
	}
	if s.Mismatched != 1 {
		t.Errorf("Mismatched: got %d, want 1", s.Mismatched)
	}
	if s.Missing != 1 {
		t.Errorf("Missing: got %d, want 1", s.Missing)
	}
}

func TestBuild_HealthyFalseWhenIssues(t *testing.T) {
	s := summariser.Build(makeFileResults())
	if s.Healthy {
		t.Error("expected Healthy=false when there are mismatches/missing keys")
	}
}

func TestBuild_HealthyTrueWhenClean(t *testing.T) {
	clean := []summariser.FileResult{
		{Name: "a.env", Results: []diff.Result{{Key: "X", Status: diff.Match}}},
	}
	s := summariser.Build(clean)
	if !s.Healthy {
		t.Error("expected Healthy=true for all-match results")
	}
}

func TestBuild_FilesSortedByName(t *testing.T) {
	s := summariser.Build(makeFileResults())
	if s.Files[0].Name != "production.env" {
		t.Errorf("expected production.env first, got %s", s.Files[0].Name)
	}
}

func TestWrite_TextContainsSummaryLine(t *testing.T) {
	var buf bytes.Buffer
	s := summariser.Build(makeFileResults())
	if err := summariser.Write(&buf, s, "text"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ISSUES FOUND") {
		t.Errorf("expected 'ISSUES FOUND' in text output, got:\n%s", out)
	}
}

func TestWrite_JSONIsValid(t *testing.T) {
	var buf bytes.Buffer
	s := summariser.Build(makeFileResults())
	if err := summariser.Write(&buf, s, "json"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	var decoded summariser.Summary
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.TotalFiles != 2 {
		t.Errorf("decoded TotalFiles: got %d, want 2", decoded.TotalFiles)
	}
}
