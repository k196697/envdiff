package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

func makeResults(statuses ...string) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: "KEY", Status: s}
	}
	return out
}

func TestCompute_Empty(t *testing.T) {
	s := scorer.Compute(nil)
	if s.Value != 100 {
		t.Errorf("expected 100, got %v", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
}

func TestCompute_AllMatch(t *testing.T) {
	res := makeResults(diff.StatusMatch, diff.StatusMatch, diff.StatusMatch)
	s := scorer.Compute(res)
	if s.Value != 100 {
		t.Errorf("expected 100, got %v", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected A, got %s", s.Grade)
	}
}

func TestCompute_AllMissing(t *testing.T) {
	res := makeResults(diff.StatusMissingInA, diff.StatusMissingInB)
	s := scorer.Compute(res)
	if s.Value != 0 {
		t.Errorf("expected 0, got %v", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected F, got %s", s.Grade)
	}
}

func TestCompute_Mixed(t *testing.T) {
	// 2 match, 1 mismatch, 1 missing → earned = 2 + 0.5 + 0 = 2.5 / 4 = 62.5
	res := makeResults(
		diff.StatusMatch,
		diff.StatusMatch,
		diff.StatusMismatch,
		diff.StatusMissingInB,
	)
	s := scorer.Compute(res)
	if s.Value != 62.5 {
		t.Errorf("expected 62.5, got %v", s.Value)
	}
	if s.Grade != "C" {
		t.Errorf("expected C, got %s", s.Grade)
	}
}

func TestCompute_Grades(t *testing.T) {
	cases := []struct {
		statuses []string
		want     string
	}{
		{makeStatuses(9, diff.StatusMatch, 1, diff.StatusMissingInA), "A"},  // 90
		{makeStatuses(3, diff.StatusMatch, 1, diff.StatusMissingInA), "C"},  // 75
		{makeStatuses(2, diff.StatusMatch, 2, diff.StatusMissingInA), "D"},  // 50
	}
	for _, tc := range cases {
		res := makeResults(tc.statuses...)
		s := scorer.Compute(res)
		if s.Grade != tc.want {
			t.Errorf("statuses %v: expected grade %s, got %s (value=%.2f)", tc.statuses, tc.want, s.Grade, s.Value)
		}
	}
}

func makeStatuses(n1 int, s1 string, n2 int, s2 string) []string {
	out := make([]string, 0, n1+n2)
	for i := 0; i < n1; i++ {
		out = append(out, s1)
	}
	for i := 0; i < n2; i++ {
		out = append(out, s2)
	}
	return out
}
