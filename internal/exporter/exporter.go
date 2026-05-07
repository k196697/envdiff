// Package exporter writes diff results to various output formats on disk.
package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/stats"
)

// Format represents a supported export format.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Options configures the export behaviour.
type Options struct {
	OutputPath string
	Format     Format
	Stats      bool
}

// Export writes diff results to the file specified in opts.
func Export(results []diff.Result, st stats.Stats, opts Options) error {
	if opts.OutputPath == "" {
		return fmt.Errorf("exporter: output path must not be empty")
	}

	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), 0o755); err != nil {
		return fmt.Errorf("exporter: create directories: %w", err)
	}

	f, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("exporter: create file: %w", err)
	}
	defer f.Close()

	switch opts.Format {
	case FormatJSON:
		return writeJSON(f, results, st, opts.Stats)
	case FormatMarkdown:
		return writeMarkdown(f, results, st, opts.Stats)
	default:
		return writeText(f, results, st, opts.Stats)
	}
}

func writeJSON(f *os.File, results []diff.Result, st stats.Stats, inclStats bool) error {
	payload := map[string]any{"results": results}
	if inclStats {
		payload["stats"] = st
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func writeText(f *os.File, results []diff.Result, st stats.Stats, inclStats bool) error {
	for _, r := range results {
		fmt.Fprintf(f, "[%s] %s (%s)\n", r.Status, r.Key, r.File)
	}
	if inclStats {
		fmt.Fprintf(f, "\nTotal: %d | Match: %d | Mismatch: %d | Missing: %d\n",
			st.Total, st.Match, st.Mismatch, st.Missing)
	}
	return nil
}

func writeMarkdown(f *os.File, results []diff.Result, st stats.Stats, inclStats bool) error {
	fmt.Fprintln(f, "# envdiff Report\n")
	fmt.Fprintln(f, "| Key | File | Status |")
	fmt.Fprintln(f, "|-----|------|--------|")
	for _, r := range results {
		fmt.Fprintf(f, "| %s | %s | %s |\n",
			escapeMD(r.Key), escapeMD(r.File), r.Status)
	}
	if inclStats {
		fmt.Fprintf(f, "\n**Total:** %d | **Match:** %d | **Mismatch:** %d | **Missing:** %d\n",
			st.Total, st.Match, st.Mismatch, st.Missing)
	}
	return nil
}

func escapeMD(s string) string {
	return strings.NewReplacer("|", "\\|", "`", "\\`").Replace(s)
}
