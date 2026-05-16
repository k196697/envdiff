// Package envnormaliser standardises env var keys and values across loaded
// environments — trimming whitespace, normalising key casing, and collapsing
// duplicate whitespace in values.
package envnormaliser

import (
	"strings"
	"unicode"
)

// Options controls normalisation behaviour.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// CollapseValues replaces runs of internal whitespace in values with a
	// single space.
	CollapseValues bool
}

// DefaultOptions returns a sensible default set of options.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimValues:     true,
		CollapseValues: false,
	}
}

// Apply normalises a copy of env according to opts. The original map is never
// mutated.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		nk := normaliseKey(k, opts)
		nv := normaliseValue(v, opts)
		out[nk] = nv
	}
	return out
}

// ApplyAll normalises every environment in the supplied map, returning a new
// map with the same names but normalised contents.
func ApplyAll(envs map[string]map[string]string, opts Options) map[string]map[string]string {
	out := make(map[string]map[string]string, len(envs))
	for name, env := range envs {
		out[name] = Apply(env, opts)
	}
	return out
}

func normaliseKey(k string, opts Options) string {
	k = strings.TrimSpace(k)
	if opts.UppercaseKeys {
		k = strings.ToUpper(k)
	}
	return k
}

func normaliseValue(v string, opts Options) string {
	if opts.TrimValues {
		v = strings.TrimSpace(v)
	}
	if opts.CollapseValues {
		v = collapseSpaces(v)
	}
	return v
}

func collapseSpaces(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return b.String()
}
