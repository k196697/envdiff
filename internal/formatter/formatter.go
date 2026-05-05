package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// jsonResult is the structure used for JSON output.
type jsonResult struct {
	MissingInA  []string            `json:"missing_in_a,omitempty"`
	MissingInB  []string            `json:"missing_in_b,omitempty"`
	Mismatched  map[string][2]string `json:"mismatched,omitempty"`
	MatchCount  int                 `json:"match_count"`
}

// Write writes the diff results to w in the specified format.
func Write(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, results)
	default:
		return writeText(w, results)
	}
}

func writeText(w io.Writer, results []diff.Result) error {
	for _, r := range results {
		switch r.Status {
		case diff.StatusMatch:
			fmt.Fprintf(w, "  [=] %s\n", r.Key)
		case diff.StatusMismatch:
			fmt.Fprintf(w, "  [~] %s: %q != %q\n", r.Key, r.ValueA, r.ValueB)
		case diff.StatusMissingInA:
			fmt.Fprintf(w, "  [+] %s (only in B)\n", r.Key)
		case diff.StatusMissingInB:
			fmt.Fprintf(w, "  [-] %s (only in A)\n", r.Key)
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []diff.Result) error {
	out := jsonResult{
		Mismatched: make(map[string][2]string),
	}
	for _, r := range results {
		switch r.Status {
		case diff.StatusMatch:
			out.MatchCount++
		case diff.StatusMismatch:
			out.Mismatched[r.Key] = [2]string{r.ValueA, r.ValueB}
		case diff.StatusMissingInA:
			out.MissingInA = append(out.MissingInA, r.Key)
		case diff.StatusMissingInB:
			out.MissingInB = append(out.MissingInB, r.Key)
		}
	}
	sort.Strings(out.MissingInA)
	sort.Strings(out.MissingInB)
	if len(out.Mismatched) == 0 {
		out.Mismatched = nil
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
