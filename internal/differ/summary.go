package differ

import (
	"github.com/user/envdiff/internal/diff"
)

// PairSummary holds aggregated counts for a single baseline→target comparison.
type PairSummary struct {
	Baseline  string
	Target    string
	Match     int
	Mismatch  int
	MissingA  int
	MissingB  int
	Total     int
	Clean     bool
}

// Summarise converts a slice of FileResult into per-pair summaries.
func Summarise(results []FileResult) []PairSummary {
	out := make([]PairSummary, 0, len(results))
	for _, fr := range results {
		ps := PairSummary{
			Baseline: fr.Baseline,
			Target:   fr.Target,
			Total:    len(fr.Results),
		}
		for _, r := range fr.Results {
			switch r.Status {
			case diff.StatusMatch:
				ps.Match++
			case diff.StatusMismatch:
				ps.Mismatch++
			case diff.StatusMissingA:
				ps.MissingA++
			case diff.StatusMissingB:
				ps.MissingB++
			}
		}
		ps.Clean = ps.Mismatch == 0 && ps.MissingA == 0 && ps.MissingB == 0
		out = append(out, ps)
	}
	return out
}
