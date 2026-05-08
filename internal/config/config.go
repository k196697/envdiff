// Package config parses CLI flags and environment variables into a Config struct.
package config

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds all runtime options for envdiff.
type Config struct {
	Files      []string
	Dir        string
	Format     string
	Exclude    []string
	Prefix     string
	SortBy     string
	OutputPath string
	ExportFmt  string
	ExportStats bool
}

var validFormats = map[string]bool{"text": true, "json": true}
var validExportFmts = map[string]bool{"text": true, "json": true, "markdown": true}
var validSortBy = map[string]bool{"key": true, "status": true, "file": true}

// Parse reads os.Args using the provided FlagSet and returns a Config.
func Parse(fs *flag.FlagSet, args []string) (*Config, error) {
	var (
		dir         string
		format      string
		exclude     string
		prefix      string
		sortBy      string
		outputPath  string
		exportFmt   string
		exportStats bool
	)

	fs.StringVar(&dir, "dir", "", "directory containing .env files")
	fs.StringVar(&format, "format", "text", "output format: text|json")
	fs.StringVar(&exclude, "exclude", "", "comma-separated keys to exclude")
	fs.StringVar(&prefix, "prefix", "", "only compare keys with this prefix")
	fs.StringVar(&sortBy, "sort", "key", "sort results by: key|status|file")
	fs.StringVar(&outputPath, "output", "", "write report to this file path")
	fs.StringVar(&exportFmt, "export-format", "text", "export format: text|json|markdown")
	fs.BoolVar(&exportStats, "export-stats", false, "include stats in exported report")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if !validFormats[format] {
		return nil, fmt.Errorf("invalid format %q: must be text or json", format)
	}
	if !validExportFmts[exportFmt] {
		return nil, fmt.Errorf("invalid export-format %q: must be text, json or markdown", exportFmt)
	}
	if !validSortBy[sortBy] {
		return nil, fmt.Errorf("invalid sort %q: must be key, status or file", sortBy)
	}

	// Require at least one input source: explicit files or a directory.
	if dir == "" && len(fs.Args()) == 0 {
		return nil, fmt.Errorf("no input provided: specify files as arguments or use -dir")
	}

	var excludeKeys []string
	if exclude != "" {
		for _, k := range strings.Split(exclude, ",") {
			if t := strings.TrimSpace(k); t != "" {
				excludeKeys = append(excludeKeys, t)
			}
		}
	}

	return &Config{
		Files:       fs.Args(),
		Dir:         dir,
		Format:      format,
		Exclude:     excludeKeys,
		Prefix:      prefix,
		SortBy:      sortBy,
		OutputPath:  outputPath,
		ExportFmt:   exportFmt,
		ExportStats: exportStats,
	}, nil
}
