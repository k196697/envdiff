// Package trimmer detects and reports keys whose values contain leading
// or trailing whitespace, which can cause subtle runtime mismatches.
package trimmer

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Issue describes a single whitespace anomaly found in a diff result.
type Issue struct {
	Key      string `json:"key"`
	File     string `json:"file"`
	Value    string `json:"value"`
	Leading  bool   `json:"leading"`
	Trailing bool   `json:"trailing"`
}

// Detect scans diff results for values with leading or trailing whitespace
// and returns a slice of Issues ordered by file then key.
func Detect(results []diff.Result) []Issue {
	var issues []Issue

	for _, r := range results {
		for file, val := range r.Values {
			leading := val != "" && val != strings.TrimLeft(val, " \t")
			trailing := val != "" && val != strings.TrimRight(val, " \t")
			if leading || trailing {
				issues = append(issues, Issue{
					Key:      r.Key,
					File:     file,
					Value:    val,
					Leading:  leading,
					Trailing: trailing,
				})
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File
		}
		return issues[i].Key < issues[j].Key
	})

	return issues
}

// WriteReport writes the detected issues to w in the requested format
// ("text" or "json"). An unknown format falls back to text.
func WriteReport(w io.Writer, issues []Issue, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return writeJSON(w, issues)
	default:
		return writeText(w, issues)
	}
}

func writeText(w io.Writer, issues []Issue) error {
	if len(issues) == 0 {
		_, err := fmt.Fprintln(w, "trimmer: no whitespace issues found")
		return err
	}
	for _, iss := range issues {
		kind := describeKind(iss)
		if _, err := fmt.Fprintf(w, "[TRIM] %s (%s): %s whitespace in value %q\n",
			iss.Key, iss.File, kind, iss.Value); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, issues []Issue) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(issues)
}

func describeKind(iss Issue) string {
	switch {
	case iss.Leading && iss.Trailing:
		return "leading+trailing"
	case iss.Leading:
		return "leading"
	default:
		return "trailing"
	}
}
