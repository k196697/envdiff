// Package renamer provides utilities for detecting and suggesting key renames
// across .env files by identifying keys that are present in one environment
// but absent in another with similar values.
package renamer

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion represents a potential rename from one key name to another.
type Suggestion struct {
	From  string
	To    string
	Score float64 // 0.0–1.0 similarity score
}

// Options controls how rename detection behaves.
type Options struct {
	MinScore float64 // minimum similarity score to include (default 0.6)
}

var defaults = Options{
	MinScore: 0.6,
}

// Detect analyses diff results and returns rename suggestions for keys that
// appear missing in one file but whose values closely match an extra key in
// another file.
func Detect(results []diff.Result, opts *Options) []Suggestion {
	if opts == nil {
		opts = &defaults
	}

	missingInB := []diff.Result{}
	missingInA := []diff.Result{}

	for _, r := range results {
		switch r.Status {
		case diff.MissingInB:
			missingInB = append(missingInB, r)
		case diff.MissingInA:
			missingInA = append(missingInA, r)
		}
	}

	var suggestions []Suggestion

	for _, a := range missingInB {
		for _, b := range missingInA {
			score := similarity(a.Key, b.Key)
			if score >= opts.MinScore {
				suggestions = append(suggestions, Suggestion{
					From:  a.Key,
					To:    b.Key,
					Score: score,
				})
			}
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Score != suggestions[j].Score {
			return suggestions[i].Score > suggestions[j].Score
		}
		return suggestions[i].From < suggestions[j].From
	})

	return suggestions
}

// similarity returns a normalised score between two strings using a simple
// token-overlap heuristic after splitting on underscores.
func similarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	aTokens := tokenise(a)
	bTokens := tokenise(b)

	shared := 0
	for _, t := range aTokens {
		for _, u := range bTokens {
			if t == u {
				shared++
				break
			}
		}
	}

	total := len(aTokens) + len(bTokens)
	if total == 0 {
		return 0
	}
	return float64(2*shared) / float64(total)
}

func tokenise(key string) []string {
	return strings.Split(strings.ToLower(key), "_")
}
