// Package scorer computes a numeric health score for a set of diff results.
// The score ranges from 0 (all keys missing or mismatched) to 100 (all match).
package scorer

import (
	"math"

	"github.com/user/envdiff/internal/diff"
)

// Weights applied to each result status when computing the score.
const (
	WeightMatch    = 1.0
	WeightMismatch = 0.5
	WeightMissing  = 0.0
)

// Score holds the computed health score and a human-readable grade.
type Score struct {
	Value float64 // 0–100
	Grade string  // A, B, C, D, F
}

// Compute calculates a health score from a slice of diff results.
// Returns a Score with Value == 100 and Grade == "A" when results is empty.
func Compute(results []diff.Result) Score {
	if len(results) == 0 {
		return Score{Value: 100, Grade: "A"}
	}

	var earned, total float64
	for _, r := range results {
		total += WeightMatch
		switch r.Status {
		case diff.StatusMatch:
			earned += WeightMatch
		case diff.StatusMismatch:
			earned += WeightMismatch
		case diff.StatusMissingInA, diff.StatusMissingInB:
			earned += WeightMissing
		}
	}

	value := math.Round((earned/total)*100*100) / 100
	return Score{
		Value: value,
		Grade: grade(value),
	}
}

func grade(v float64) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}
