// Package auditor provides functionality for generating an audit trail
// of differences found between .env files, recording what changed,
// when the audit was run, and which files were compared.
package auditor

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Entry represents a single audit record for one key across compared files.
type Entry struct {
	Key    string
	Status string
	Files  []string
}

// Report holds the full audit output for a comparison run.
type Report struct {
	RunAt   time.Time
	Files   []string
	Entries []Entry
	Total   int
	Issues  int
}

// Build constructs an audit Report from diff results and the list of
// compared file paths. t is the timestamp to record; pass time.Now()
// in production callers.
func Build(results []diff.Result, files []string, t time.Time) Report {
	entries := make([]Entry, 0, len(results))
	issues := 0

	for _, r := range results {
		e := Entry{
			Key:    r.Key,
			Status: string(r.Status),
			Files:  files,
		}
		entries = append(entries, e)
		if r.Status != diff.StatusMatch {
			issues++
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Report{
		RunAt:   t,
		Files:   files,
		Entries: entries,
		Total:   len(entries),
		Issues:  issues,
	}
}

// Write formats the audit Report as human-readable text and writes it to w.
func Write(r Report, w io.Writer) error {
	_, err := fmt.Fprintf(w, "Audit Report\nRun at : %s\nFiles  : %v\nTotal  : %d\nIssues : %d\n\n",
		r.RunAt.Format(time.RFC3339), r.Files, r.Total, r.Issues)
	if err != nil {
		return err
	}

	for _, e := range r.Entries {
		_, err = fmt.Fprintf(w, "  [%-10s] %s\n", e.Status, e.Key)
		if err != nil {
			return err
		}
	}
	return nil
}
