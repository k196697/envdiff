// Package redactor masks sensitive values in diff results before display or export.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DefaultSensitivePatterns are key substrings that trigger redaction.
var DefaultSensitivePatterns = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
	"passwd",
	"auth",
}

const redactedValue = "***REDACTED***"

// Options controls redactor behaviour.
type Options struct {
	// Patterns is the list of case-insensitive key substrings that trigger redaction.
	// When nil, DefaultSensitivePatterns is used.
	Patterns []string
}

// isSensitive reports whether key matches any of the provided patterns.
func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Apply returns a copy of results with sensitive values replaced by the
// redacted placeholder. Original results are not modified.
func Apply(results []diff.Result, opts Options) []diff.Result {
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}

	out := make([]diff.Result, len(results))
	for i, r := range results {
		if isSensitive(r.Key, patterns) {
			r.ValueA = redactedValue
			r.ValueB = redactedValue
		}
		out[i] = r
	}
	return out
}
