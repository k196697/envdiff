// Package patcher writes missing keys from a reference env into a target env file.
package patcher

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls patching behaviour.
type Options struct {
	// DryRun prints what would be written without modifying the file.
	DryRun bool
	// Placeholder is the value written for missing keys. Defaults to "".
	Placeholder string
}

// Result describes the outcome of a patch operation.
type Result struct {
	File        string
	KeysAdded   []string
	DryRun      bool
}

// Apply appends keys that are missing in the target file, using values from
// the reference map when available. It operates on the diff results produced
// by diff.Compare so callers control what "missing" means.
func Apply(targetPath string, results []diff.Result, opts Options) (Result, error) {
	if opts.Placeholder == "" {
		opts.Placeholder = ""
	}

	var missing []diff.Result
	for _, r := range results {
		if r.Status == diff.Missing {
			missing = append(missing, r)
		}
	}

	if len(missing) == 0 {
		return Result{File: targetPath}, nil
	}

	sort.Slice(missing, func(i, j int) bool {
		return missing[i].Key < missing[j].Key
	})

	var lines []string
	for _, r := range missing {
		val := opts.Placeholder
		if r.ValueA != "" {
			val = r.ValueA
		}
		lines = append(lines, fmt.Sprintf("%s=%s", r.Key, val))
	}

	block := "\n# patched by envdiff\n" + strings.Join(lines, "\n") + "\n"

	var added []string
	for _, r := range missing {
		added = append(added, r.Key)
	}

	if opts.DryRun {
		return Result{File: targetPath, KeysAdded: added, DryRun: true}, nil
	}

	f, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return Result{}, fmt.Errorf("patcher: open %s: %w", targetPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(block); err != nil {
		return Result{}, fmt.Errorf("patcher: write %s: %w", targetPath, err)
	}

	return Result{File: targetPath, KeysAdded: added}, nil
}
