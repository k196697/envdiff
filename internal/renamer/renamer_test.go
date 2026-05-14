package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/renamer"
)

func makeResults(pairs [][3]string) []diff.Result {
	out := make([]diff.Result, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, diff.Result{Key: p[0], Status: diff.Status(p[1]), ValueA: p[2]})
	}
	return out
}

func TestDetect_NoMissing(t *testing.T) {
	results := makeResults([][3]string{
		{"DB_HOST", string(diff.Match), "localhost"},
	})
	got := renamer.Detect(results, nil)
	if len(got) != 0 {
		t.Fatalf("expected no suggestions, got %d", len(got))
	}
}

func TestDetect_ObviousRename(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.MissingInB},
		{Key: "DB_HOSTNAME", Status: diff.MissingInA},
	}
	got := renamer.Detect(results, nil)
	if len(got) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	if got[0].From != "DB_HOST" || got[0].To != "DB_HOSTNAME" {
		t.Errorf("unexpected suggestion: %+v", got[0])
	}
	if got[0].Score <= 0 || got[0].Score > 1.0 {
		t.Errorf("score out of range: %f", got[0].Score)
	}
}

func TestDetect_BelowMinScore(t *testing.T) {
	results := []diff.Result{
		{Key: "ALPHA", Status: diff.MissingInB},
		{Key: "ZZZZ", Status: diff.MissingInA},
	}
	got := renamer.Detect(results, &renamer.Options{MinScore: 0.9})
	if len(got) != 0 {
		t.Fatalf("expected no suggestions above threshold, got %d", len(got))
	}
}

func TestDetect_SortedByScoreDesc(t *testing.T) {
	results := []diff.Result{
		{Key: "APP_SECRET_KEY", Status: diff.MissingInB},
		{Key: "DB_HOST", Status: diff.MissingInB},
		{Key: "APP_SECRET", Status: diff.MissingInA},
		{Key: "DB_HOSTNAME", Status: diff.MissingInA},
	}
	got := renamer.Detect(results, &renamer.Options{MinScore: 0.3})
	for i := 1; i < len(got); i++ {
		if got[i].Score > got[i-1].Score {
			t.Errorf("suggestions not sorted by score desc at index %d", i)
		}
	}
}

func TestDetect_ExactMatch_ScoreOne(t *testing.T) {
	score := renamer.Similarity("FOO_BAR", "FOO_BAR")
	if score != 1.0 {
		t.Errorf("expected 1.0 for identical keys, got %f", score)
	}
}

func TestDetect_CustomMinScore(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.MissingInB},
		{Key: "DB_HOSTNAME", Status: diff.MissingInA},
	}
	got := renamer.Detect(results, &renamer.Options{MinScore: 0.99})
	if len(got) != 0 {
		t.Fatalf("expected no results with very high threshold, got %d", len(got))
	}
}
