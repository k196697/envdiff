// Package duplicator detects duplicate values across keys in one or more
// parsed .env maps. Two keys are considered duplicates when they share an
// identical non-empty value, which often indicates copy-paste mistakes or
// unintentional aliasing.
package duplicator

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Entry describes a single duplicate-value group found inside one env file.
type Entry struct {
	// Value is the shared value.
	Value string
	// Keys are the keys that all carry Value, sorted alphabetically.
	Keys []string
	// File is the label / path of the env map that was inspected.
	File string
}

// Detect scans each named env map and returns every group of keys that share
// the same non-empty value. Results are sorted by file then by value.
func Detect(envs map[string]map[string]string) []Entry {
	var results []Entry

	// Collect file names so output is deterministic.
	files := make([]string, 0, len(envs))
	for f := range envs {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		env := envs[file]
		// Build an inverse map: value -> []keys
		inverse := make(map[string][]string)
		for k, v := range env {
			if v == "" {
				continue
			}
			inverse[v] = append(inverse[v], k)
		}

		// Collect values that appear more than once.
		values := make([]string, 0, len(inverse))
		for v, keys := range inverse {
			if len(keys) > 1 {
				values = append(values, v)
			}
		}
		sort.Strings(values)

		for _, v := range values {
			keys := inverse[v]
			sort.Strings(keys)
			results = append(results, Entry{
				Value: v,
				Keys:  keys,
				File:  file,
			})
		}
	}

	return results
}

// Write renders the duplicate entries to w in a human-readable format.
// If entries is empty a clean message is printed instead.
func Write(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no duplicate values detected")
		return
	}
	for _, e := range entries {
		fmt.Fprintf(w, "[%s] value %q shared by: %s\n",
			e.File, e.Value, strings.Join(e.Keys, ", "))
	}
}
