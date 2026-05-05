package loader

import (
	"fmt"
	"strings"
)

// ValidationError holds all issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Issues, "; "))
}

// Validate checks a slice of EnvFile for common issues:
//   - duplicate file paths
//   - empty parsed environments (files with no keys)
//
// Returns a *ValidationError if any issues are found, or nil if all files are valid.
func Validate(files []EnvFile) error {
	var issues []string
	seen := make(map[string]bool)

	for _, f := range files {
		if seen[f.Path] {
			issues = append(issues, fmt.Sprintf("duplicate file path: %q", f.Path))
		}
		seen[f.Path] = true

		if len(f.Env) == 0 {
			issues = append(issues, fmt.Sprintf("file %q contains no keys", f.Name))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
