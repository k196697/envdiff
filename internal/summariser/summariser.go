// Package summariser produces a human-readable or JSON summary of a full
// envdiff run, aggregating per-file diff results into a single report.
package summariser

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// FileResult holds the diff results for a single file pair.
type FileResult struct {
	Name    string
	Results []diff.Result
}

// Summary is the top-level aggregated report.
type Summary struct {
	GeneratedAt  time.Time      `json:"generated_at"`
	TotalFiles   int            `json:"total_files"`
	TotalKeys    int            `json:"total_keys"`
	Matched      int            `json:"matched"`
	Mismatched   int            `json:"mismatched"`
	Missing      int            `json:"missing"`
	Healthy      bool           `json:"healthy"`
	Files        []FileSummary  `json:"files"`
}

// FileSummary is the per-file breakdown inside a Summary.
type FileSummary struct {
	Name       string `json:"name"`
	Keys       int    `json:"keys"`
	Matched    int    `json:"matched"`
	Mismatched int    `json:"mismatched"`
	Missing    int    `json:"missing"`
}

// Build aggregates a slice of FileResult values into a Summary.
func Build(files []FileResult) Summary {
	s := Summary{
		GeneratedAt: time.Now().UTC(),
		TotalFiles:  len(files),
	}

	for _, f := range files {
		fs := FileSummary{Name: f.Name, Keys: len(f.Results)}
		for _, r := range f.Results {
			switch r.Status {
			case diff.Match:
				fs.Matched++
			case diff.Mismatch:
				fs.Mismatched++
			case diff.MissingInA, diff.MissingInB:
				fs.Missing++
			}
		}
		s.TotalKeys += fs.Keys
		s.Matched += fs.Matched
		s.Mismatched += fs.Mismatched
		s.Missing += fs.Missing
		s.Files = append(s.Files, fs)
	}

	sort.Slice(s.Files, func(i, j int) bool {
		return s.Files[i].Name < s.Files[j].Name
	})

	s.Healthy = s.Mismatched == 0 && s.Missing == 0
	return s
}

// Write renders the Summary to w in the requested format ("text" or "json").
func Write(w io.Writer, s Summary, format string) error {
	if format == "json" {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(s)
	}
	return writeText(w, s)
}

func writeText(w io.Writer, s Summary) error {
	status := "OK"
	if !s.Healthy {
		status = "ISSUES FOUND"
	}
	fmt.Fprintf(w, "envdiff summary — %s\n", status)
	fmt.Fprintf(w, "Generated : %s\n", s.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Files     : %d\n", s.TotalFiles)
	fmt.Fprintf(w, "Keys      : %d  matched=%d  mismatched=%d  missing=%d\n\n",
		s.TotalKeys, s.Matched, s.Mismatched, s.Missing)
	for _, f := range s.Files {
		fmt.Fprintf(w, "  %-40s keys=%-4d matched=%-4d mismatched=%-4d missing=%d\n",
			f.Name, f.Keys, f.Matched, f.Mismatched, f.Missing)
	}
	return nil
}
