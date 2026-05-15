// Package resolver expands variable references within .env values.
// It resolves expressions like ${OTHER_KEY} or $OTHER_KEY using the
// values already present in the same environment map.
package resolver

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the outcome of resolving a single key.
type Result struct {
	Key      string
	Original string
	Resolved string
	Refs     []string // keys that were referenced
	Missing  []string // referenced keys that had no value
}

// Resolve expands all variable references in env and returns one Result
// per key. Keys are processed in a single pass; forward references are
// supported as long as there are no cycles.
func Resolve(env map[string]string) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		resolved, refs, missing := expand(v, env)
		results = append(results, Result{
			Key:      k,
			Original: v,
			Resolved: resolved,
			Refs:     refs,
			Missing:  missing,
		})
	}
	return results
}

// expand performs the substitution for a single value string.
func expand(value string, env map[string]string) (resolved string, refs []string, missing []string) {
	seen := map[string]bool{}
	resolved = refPattern.ReplaceAllStringFunc(value, func(match string) string {
		key := extractKey(match)
		if !seen[key] {
			seen[key] = true
			refs = append(refs, key)
		}
		if val, ok := env[key]; ok {
			return val
		}
		missing = append(missing, key)
		return fmt.Sprintf("<MISSING:%s>", key)
	})
	return resolved, refs, missing
}

// extractKey pulls the variable name out of a ${VAR} or $VAR match.
func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
