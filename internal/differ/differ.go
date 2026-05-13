// Package differ provides multi-file diff capabilities, comparing
// all loaded environments against a designated baseline file.
package differ

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// FileResult holds the diff results for one file compared to the baseline.
type FileResult struct {
	Baseline string
	Target   string
	Results  []diff.Result
}

// Options controls how the multi-file diff is performed.
type Options struct {
	// BaselineIndex is the index of the file to treat as the baseline (default 0).
	BaselineIndex int
}

// RunAll compares every file in envs against the baseline file.
// envs is a map of label -> parsed env map.
func RunAll(envs map[string]map[string]string, baselineLabel string) ([]FileResult, error) {
	baseline, ok := envs[baselineLabel]
	if !ok {
		return nil, fmt.Errorf("differ: baseline label %q not found in envs", baselineLabel)
	}

	var out []FileResult
	for label, env := range envs {
		if label == baselineLabel {
			continue
		}
		results := diff.Compare(baseline, env)
		out = append(out, FileResult{
			Baseline: baselineLabel,
			Target:   label,
			Results:  results,
		})
	}
	return out, nil
}

// RunFiles loads and compares all provided file paths against the first path
// as the baseline, returning one FileResult per non-baseline file.
func RunFiles(paths []string) ([]FileResult, error) {
	if len(paths) < 2 {
		return nil, fmt.Errorf("differ: at least two files required, got %d", len(paths))
	}

	baseline, err := parser.ParseFile(paths[0])
	if err != nil {
		return nil, fmt.Errorf("differ: parsing baseline %s: %w", paths[0], err)
	}

	var out []FileResult
	for _, p := range paths[1:] {
		env, err := parser.ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("differ: parsing %s: %w", p, err)
		}
		results := diff.Compare(baseline, env)
		out = append(out, FileResult{
			Baseline: paths[0],
			Target:   p,
			Results:  results,
		})
	}
	return out, nil
}
