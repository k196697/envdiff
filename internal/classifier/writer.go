package classifier

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes classified results to w in the requested format.
// Supported formats: "text" (default), "json".
func WriteReport(w io.Writer, results []Result, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return writeJSON(w, results)
	default:
		return writeText(w, results)
	}
}

func writeText(w io.Writer, results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no results to classify")
		return err
	}

	// sort: critical first, then warning, then info; ties broken by key
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		ri, rj := severityOrder(sorted[i].Severity), severityOrder(sorted[j].Severity)
		if ri != rj {
			return ri < rj
		}
		return sorted[i].Diff.Key < sorted[j].Diff.Key
	})

	for _, r := range sorted {
		_, err := fmt.Fprintf(w, "[%-8s] %-30s %s\n",
			strings.ToUpper(string(r.Severity)),
			r.Diff.Key,
			r.Diff.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []Result) error {
	type row struct {
		Key      string `json:"key"`
		Status   string `json:"status"`
		Severity string `json:"severity"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{
			Key:      r.Diff.Key,
			Status:   string(r.Diff.Status),
			Severity: string(r.Severity),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func severityOrder(s Severity) int {
	switch s {
	case SeverityCritical:
		return 0
	case SeverityWarning:
		return 1
	default:
		return 2
	}
}
