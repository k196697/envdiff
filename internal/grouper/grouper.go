// Package grouper organises diff results by key prefix, allowing callers
// to see which logical sections of an environment have issues.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Group holds all results that share a common prefix.
type Group struct {
	Prefix  string
	Results []diff.Result
}

// Options controls how grouping is performed.
type Options struct {
	// Separator is the character used to split key prefixes (default "_").
	Separator string
	// Depth is how many separator-delimited segments form the prefix (default 1).
	Depth int
}

func defaults(o Options) Options {
	if o.Separator == "" {
		o.Separator = "_"
	}
	if o.Depth < 1 {
		o.Depth = 1
	}
	return o
}

// ByPrefix groups results by the leading segments of each key.
// Results whose key has no separator are placed under the group "(other)".
func ByPrefix(results []diff.Result, opts Options) []Group {
	opts = defaults(opts)

	buckets := make(map[string][]diff.Result)
	for _, r := range results {
		prefix := extractPrefix(r.Key, opts.Separator, opts.Depth)
		buckets[prefix] = append(buckets[prefix], r)
	}

	groups := make([]Group, 0, len(buckets))
	for prefix, res := range buckets {
		groups = append(groups, Group{Prefix: prefix, Results: res})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})
	return groups
}

func extractPrefix(key, sep string, depth int) string {
	parts := strings.Split(key, sep)
	if len(parts) <= depth {
		return "(other)"
	}
	return strings.Join(parts[:depth], sep)
}
