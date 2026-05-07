// Package templater generates a template .env file from a set of diff results,
// producing a file with all known keys and placeholder values for missing ones.
package templater

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls template generation behaviour.
type Options struct {
	// Placeholder is the value written for keys that are missing or mismatched.
	// Defaults to "REPLACE_ME" when empty.
	Placeholder string
	// IncludeValues keeps the original value for keys that match across all files.
	IncludeValues bool
}

// Generate writes a template .env file to w based on the provided diff results.
// Keys are sorted alphabetically. Keys that are present and matching retain
// their value when IncludeValues is true; all other keys receive the placeholder.
func Generate(w io.Writer, results []diff.Result, opts Options) error {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "REPLACE_ME"
	}

	// Deduplicate keys, keeping the best available value.
	type entry struct {
		value  string
		keepValue bool
	}
	seen := make(map[string]entry)

	for _, r := range results {
		key := r.Key
		if _, exists := seen[key]; exists {
			continue
		}
		switch r.Status {
		case diff.Match:
			seen[key] = entry{value: r.ValueA, keepValue: opts.IncludeValues}
		default:
			seen[key] = entry{value: placeholder, keepValue: false}
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		e := seen[k]
		v := placeholder
		if e.keepValue {
			v = e.value
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	_, err := io.WriteString(w, sb.String())
	return err
}
