// Package linter checks .env file entries for common issues such as
// empty values, keys with suspicious characters, or values that look
// like unresolved placeholders.
package linter

import (
	"fmt"
	"strings"
)

// Issue represents a single linting problem found in an env map.
type Issue struct {
	File string
	Key  string
	Msg  string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: [%s] %s", i.File, i.Key, i.Msg)
}

// Lint inspects a named env map and returns any issues found.
func Lint(name string, env map[string]string) []Issue {
	var issues []Issue

	for key, val := range env {
		if containsSpace(key) {
			issues = append(issues, Issue{
				File: name,
				Key:  key,
				Msg:  "key contains whitespace",
			})
		}

		if strings.ToUpper(key) != key {
			issues = append(issues, Issue{
				File: name,
				Key:  key,
				Msg:  "key is not uppercase",
			})
		}

		if val == "" {
			issues = append(issues, Issue{
				File: name,
				Key:  key,
				Msg:  "value is empty",
			})
		}

		if isUnresolvedPlaceholder(val) {
			issues = append(issues, Issue{
				File: name,
				Key:  key,
				Msg:  fmt.Sprintf("value looks like an unresolved placeholder: %q", val),
			})
		}
	}

	return issues
}

func containsSpace(s string) bool {
	return strings.ContainsAny(s, " \t")
}

func isUnresolvedPlaceholder(val string) bool {
	return (strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}")) ||
		(strings.HasPrefix(val, "<") && strings.HasSuffix(val, ">")) ||
		strings.EqualFold(val, "todo") ||
		strings.EqualFold(val, "changeme") ||
		strings.EqualFold(val, "fixme")
}
