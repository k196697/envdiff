// Package envprinter renders a parsed env map as formatted output,
// supporting text and JSON formats for display or piping.
package envprinter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how the env map is printed.
type Options struct {
	// Format is either "text" (default) or "json".
	Format string
	// SortKeys sorts keys alphabetically when true.
	SortKeys bool
	// MaskValues replaces all values with "***" when true.
	MaskValues bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:   "text",
		SortKeys: true,
	}
}

// Write prints the env map to w using the provided options.
func Write(w io.Writer, env map[string]string, opts Options) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeJSON(w, env, opts)
	default:
		return writeText(w, env, opts)
	}
}

func writeText(w io.Writer, env map[string]string, opts Options) error {
	keys := sortedKeys(env, opts.SortKeys)
	for _, k := range keys {
		v := env[k]
		if opts.MaskValues && v != "" {
			v = "***"
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, env map[string]string, opts Options) error {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.MaskValues && v != "" {
			out[k] = "***"
		} else {
			out[k] = v
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedKeys(env map[string]string, doSort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if doSort {
		sort.Strings(keys)
	}
	return keys
}
