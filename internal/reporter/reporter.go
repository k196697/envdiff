package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to the given writer.
func Report(results []diff.Result, fileA, fileB string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "Comparing: %s <-> %s\n", fileA, fileB)
	fmt.Fprintln(w, strings.Repeat("-", 50))

	// Sort results by key for deterministic output
	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		switch r.Status {
		case diff.StatusMatch:
			fmt.Fprintf(w, "  [OK]      %s\n", r.Key)
		case diff.StatusMismatch:
			fmt.Fprintf(w, "  [MISMATCH] %s\n", r.Key)
			fmt.Fprintf(w, "             %s: %q\n", fileA, r.ValueA)
			fmt.Fprintf(w, "             %s: %q\n", fileB, r.ValueB)
		case diff.StatusMissingInB:
			fmt.Fprintf(w, "  [MISSING]  %s (not in %s)\n", r.Key, fileB)
		case diff.StatusMissingInA:
			fmt.Fprintf(w, "  [MISSING]  %s (not in %s)\n", r.Key, fileA)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))
	summary := diff.Summary(results)
	fmt.Fprintf(w, "Summary: %d match, %d mismatch, %d missing\n",
		summary["match"], summary["mismatch"], summary["missing"])
}
