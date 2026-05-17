package envchain

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteReport writes the resolved chain results to w in the given format
// ("text" or "json"). Unrecognised formats fall back to text.
func WriteReport(w io.Writer, results []Result, format string) error {
	switch format {
	case "json":
		return writeJSON(w, results)
	default:
		return writeText(w, results)
	}
}

func writeText(w io.Writer, results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no keys resolved")
		return err
	}
	for _, r := range results {
		line := fmt.Sprintf("%-30s = %s  (from: %s)", r.Key, r.Value, r.Source)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
		for _, ov := range r.Overridden {
			if _, err := fmt.Fprintf(w, "  overrides: %s = %s (in %s)\n", r.Key, ov.Value, ov.Source); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []Result) error {
	type jsonOverride struct {
		Source string `json:"source"`
		Value  string `json:"value"`
	}
	type jsonResult struct {
		Key        string         `json:"key"`
		Value      string         `json:"value"`
		Source     string         `json:"source"`
		Overridden []jsonOverride `json:"overridden,omitempty"`
	}
	out := make([]jsonResult, 0, len(results))
	for _, r := range results {
		jr := jsonResult{Key: r.Key, Value: r.Value, Source: r.Source}
		for _, ov := range r.Overridden {
			jr.Overridden = append(jr.Overridden, jsonOverride{Source: ov.Source, Value: ov.Value})
		}
		out = append(out, jr)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
