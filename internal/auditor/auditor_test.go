package auditor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/diff"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch},
		{Key: "API_KEY", Status: diff.StatusMismatch},
		{Key: "SECRET", Status: diff.StatusMissingInB},
	}
}

func TestBuild_CountsIssues(t *testing.T) {
	r := auditor.Build(makeResults(), []string{"a.env", "b.env"}, fixedTime)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.Issues != 2 {
		t.Errorf("expected Issues=2, got %d", r.Issues)
	}
}

func TestBuild_SortedByKey(t *testing.T) {
	r := auditor.Build(makeResults(), []string{"a.env", "b.env"}, fixedTime)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	expected := []string{"API_KEY", "DB_HOST", "SECRET"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestBuild_RecordsTimestampAndFiles(t *testing.T) {
	files := []string{"prod.env", "staging.env"}
	r := auditor.Build(makeResults(), files, fixedTime)
	if !r.RunAt.Equal(fixedTime) {
		t.Errorf("unexpected RunAt: %v", r.RunAt)
	}
	if len(r.Files) != 2 || r.Files[0] != "prod.env" {
		t.Errorf("unexpected Files: %v", r.Files)
	}
}

func TestWrite_ContainsKeyAndStatus(t *testing.T) {
	r := auditor.Build(makeResults(), []string{"a.env", "b.env"}, fixedTime)
	var buf bytes.Buffer
	if err := auditor.Write(r, &buf); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"API_KEY", "mismatch", "DB_HOST", "match", "SECRET", "missing_in_b"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	r := auditor.Build(nil, []string{"a.env"}, fixedTime)
	var buf bytes.Buffer
	if err := auditor.Write(r, &buf); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if r.Total != 0 || r.Issues != 0 {
		t.Errorf("expected zero totals for empty results")
	}
}
