// Package masker rewrites diff results so that sensitive values are replaced
// with a configurable mask string before any output is produced.
package masker

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

const DefaultMask = "***"

// Options controls masking behaviour.
type Options struct {
	// Mask is the replacement string used for sensitive values.
	// Defaults to DefaultMask when empty.
	Mask string

	// ExtraPatterns are additional sub-strings that, when found in a key
	// (case-insensitive), cause the value to be masked.
	ExtraPatterns []string
}

var defaultSensitivePatterns = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"auth", "credential", "private", "cert", "key",
}

// Apply returns a new slice of diff.Result with sensitive values replaced by
// the mask string. The original slice is never mutated.
func Apply(results []diff.Result, opts Options) []diff.Result {
	mask := opts.Mask
	if mask == "" {
		mask = DefaultMask
	}

	patterns := make([]string, len(defaultSensitivePatterns)+len(opts.ExtraPatterns))
	copy(patterns, defaultSensitivePatterns)
	copy(patterns[len(defaultSensitivePatterns):], opts.ExtraPatterns)

	out := make([]diff.Result, len(results))
	for i, r := range results {
		if isSensitive(r.Key, patterns) {
			r.ValueA = maskNonEmpty(r.ValueA, mask)
			r.ValueB = maskNonEmpty(r.ValueB, mask)
		}
		out[i] = r
	}
	return out
}

func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func maskNonEmpty(v, mask string) string {
	if v == "" {
		return v
	}
	return mask
}
