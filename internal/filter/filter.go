package filter

import "strings"

// Options holds filtering configuration for env key comparisons.
type Options struct {
	// Prefix restricts comparison to keys starting with this prefix (case-insensitive).
	Prefix string
	// Exclude is a list of exact key names to skip during comparison.
	Exclude []string
}

// Apply returns a filtered copy of the given env map based on the provided Options.
// Keys not matching the prefix or present in the exclusion list are removed.
func Apply(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))
	excludeSet := buildExcludeSet(opts.Exclude)

	for k, v := range env {
		if excludeSet[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(strings.ToUpper(k), strings.ToUpper(opts.Prefix)) {
			continue
		}
		result[k] = v
	}
	return result
}

// buildExcludeSet converts a slice of key names into a lookup map.
func buildExcludeSet(keys []string) map[string]bool {
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	return set
}
