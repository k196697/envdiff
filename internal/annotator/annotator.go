// Package annotator attaches human-readable descriptions to diff results
// by matching keys against a known annotation map or a loaded annotation file.
package annotator

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Annotation holds a key and its description.
type Annotation struct {
	Key         string
	Description string
}

// AnnotatedResult wraps a diff.Result with an optional description.
type AnnotatedResult struct {
	diff.Result
	Description string
}

// LoadFile reads an annotation file where each non-blank, non-comment line
// has the form:  KEY=Some human readable description
func LoadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	annotations := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		desc := strings.TrimSpace(line[idx+1:])
		if key != "" {
			annotations[key] = desc
		}
	}
	return annotations, scanner.Err()
}

// Apply enriches a slice of diff.Result with descriptions from the provided
// annotation map. Results without a matching annotation receive an empty
// Description field.
func Apply(results []diff.Result, annotations map[string]string) []AnnotatedResult {
	out := make([]AnnotatedResult, len(results))
	for i, r := range results {
		out[i] = AnnotatedResult{
			Result:      r,
			Description: annotations[r.Key],
		}
	}
	return out
}
