// Package promoter suggests which keys from a lower-priority environment
// should be promoted to a higher-priority environment (e.g. staging → production).
package promoter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Suggestion describes a key that is a candidate for promotion.
type Suggestion struct {
	Key       string
	FromEnv   string
	FromValue string
	ToEnv     string
	Reason    string
}

// Options controls which keys are considered for promotion.
type Options struct {
	// OnlyMissing limits suggestions to keys absent in the target env.
	OnlyMissing bool
	// IncludeMismatch also suggests keys whose values differ.
	IncludeMismatch bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{OnlyMissing: true, IncludeMismatch: false}
}

// Promote compares fromEnv against toEnv and returns promotion suggestions.
// fromEnv is the source (lower priority), toEnv is the target (higher priority).
func Promote(fromEnv, toEnv map[string]string, fromName, toName string, opts Options) []Suggestion {
	var suggestions []Suggestion

	keys := sortedKeys(fromEnv)
	for _, k := range keys {
		fromVal := fromEnv[k]
		toVal, exists := toEnv[k]

		if !exists && opts.OnlyMissing {
			suggestions = append(suggestions, Suggestion{
				Key:       k,
				FromEnv:   fromName,
				FromValue: fromVal,
				ToEnv:     toName,
				Reason:    fmt.Sprintf("key missing in %s", toName),
			})
			continue
		}

		if exists && opts.IncludeMismatch && fromVal != toVal {
			suggestions = append(suggestions, Suggestion{
				Key:       k,
				FromEnv:   fromName,
				FromValue: fromVal,
				ToEnv:     toName,
				Reason:    fmt.Sprintf("value differs between %s and %s", fromName, toName),
			})
		}
	}

	return suggestions
}

// WriteReport writes a human-readable or JSON promotion report to w.
func WriteReport(w io.Writer, suggestions []Suggestion, format string) error {
	if strings.EqualFold(format, "json") {
		return writeJSON(w, suggestions)
	}
	return writeText(w, suggestions)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
