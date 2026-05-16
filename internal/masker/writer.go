package masker

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// WriteReport writes masked results to w in the requested format.
// Supported formats: "text" (default), "json".
func WriteReport(w io.Writer, results []diff.Result, format string) error {
	switch format {
	case "json":
		return writeJSON(w, results)
	default:
		return writeText(w, results)
	}
}

func writeText(w io.Writer, results []diff.Result) error {
	for _, r := range results {
		var line string
		switch r.Status {
		case diff.Match:
			line = fmt.Sprintf("  OK  %s=%s", r.Key, r.ValueA)
		case diff.Mismatch:
			line = fmt.Sprintf(" DIFF %s: %s vs %s", r.Key, r.ValueA, r.ValueB)
		case diff.MissingInA:
			line = fmt.Sprintf(" MISS %s (missing in first file)", r.Key)
		case diff.MissingInB:
			line = fmt.Sprintf(" MISS %s (missing in second file)", r.Key)
		default:
			line = fmt.Sprintf("  ?   %s", r.Key)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []diff.Result) error {
	type row struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		ValueA string `json:"value_a"`
		ValueB string `json:"value_b"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{
			Key:    r.Key,
			Status: string(r.Status),
			ValueA: r.ValueA,
			ValueB: r.ValueB,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
