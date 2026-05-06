// Package stats computes summary statistics over a set of diff results.
package stats

import (
	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregate counts derived from a slice of diff results.
type Stats struct {
	Total    int
	Matching int
	Missing  int
	Mismatch int
	// Coverage is the percentage of keys that match (0-100).
	Coverage float64
}

// Compute calculates statistics from a map of filename -> diff results.
func Compute(results map[string][]diff.Result) Stats {
	var s Stats

	seen := make(map[string]struct{})

	for _, fileResults := range results {
		for _, r := range fileResults {
			key := r.Key
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				s.Total++
			}

			switch r.Status {
			case diff.StatusMatch:
				s.Matching++
			case diff.StatusMissing:
				s.Missing++
			case diff.StatusMismatch:
				s.Mismatch++
			}
		}
	}

	if s.Total > 0 {
		s.Coverage = float64(s.Matching) / float64(s.Total) * 100
	}

	return s
}

// IsClean returns true when there are no missing or mismatched keys.
func (s Stats) IsClean() bool {
	return s.Missing == 0 && s.Mismatch == 0
}
