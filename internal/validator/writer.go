package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteReport writes violations to w in the requested format ("text" or "json").
// An empty format defaults to "text".
func WriteReport(w io.Writer, violations []Violation, format string) error {
	if format == "" {
		format = "text"
	}
	switch format {
	case "json":
		return writeJSON(w, violations)
	default:
		return writeText(w, violations)
	}
}

func writeText(w io.Writer, violations []Violation) error {
	if len(violations) == 0 {
		_, err := fmt.Fprintln(w, "validation ok — no violations found")
		return err
	}

	sorted := make([]Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].File != sorted[j].File {
			return sorted[i].File < sorted[j].File
		}
		return sorted[i].Key < sorted[j].Key
	})

	for _, v := range sorted {
		if _, err := fmt.Fprintf(w, "[%s] %s (%s): %s\n", v.File, v.Key, v.Rule, v.Message); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, violations []Violation) error {
	type payload struct {
		Violations []Violation `json:"violations"`
		Total       int         `json:"total"`
	}
	p := payload{Violations: violations, Total: len(violations)}
	if p.Violations == nil {
		p.Violations = []Violation{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
