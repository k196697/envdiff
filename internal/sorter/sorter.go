// Package sorter provides utilities for sorting and ordering diff results
// consistently across output formats.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortOrder defines the ordering strategy for diff results.
type SortOrder int

const (
	// ByKey sorts results alphabetically by key name.
	ByKey SortOrder = iota
	// ByStatus groups results by their diff status (missing, mismatch, match).
	ByStatus
	// ByFile sorts results by the file they originate from.
	ByFile
)

// statusRank assigns a numeric priority to each diff status for ByStatus ordering.
func statusRank(r diff.Result) int {
	switch {
	case r.MissingInA || r.MissingInB:
		return 0
	case r.Mismatch:
		return 1
	default:
		return 2
	}
}

// Sort returns a new slice of Results ordered according to the given SortOrder.
// The original slice is not modified.
func Sort(results []diff.Result, order SortOrder) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	switch order {
	case ByStatus:
		sort.SliceStable(out, func(i, j int) bool {
			ri, rj := statusRank(out[i]), statusRank(out[j])
			if ri != rj {
				return ri < rj
			}
			return out[i].Key < out[j].Key
		})
	case ByFile:
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].FileA != out[j].FileA {
				return out[i].FileA < out[j].FileA
			}
			return out[i].Key < out[j].Key
		})
	default: // ByKey
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}
