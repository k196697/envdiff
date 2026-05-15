package resolver

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes a human-readable or JSON report of resolution results
// to w. format must be "text" or "json" (default: "text").
func WriteReport(w io.Writer, results []Result, format string) error {
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	switch strings.ToLower(format) {
	case "json":
		return writeJSON(w, sorted)
	default:
		return writeText(w, sorted)
	}
}

func writeText(w io.Writer, results []Result) error {
	for _, r := range results {
		if r.Original == r.Resolved && len(r.Missing) == 0 {
			continue // nothing interesting to report
		}
		if len(r.Missing) > 0 {
			_, err := fmt.Fprintf(w, "[WARN] %s: unresolved refs %v → %s\n",
				r.Key, r.Missing, r.Resolved)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "[OK]   %s: %q → %q\n",
				r.Key, r.Original, r.Resolved)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
