// Package ignorer provides support for loading and matching keys
// defined in an .envdiffignore file, similar to .gitignore semantics.
package ignorer

import (
	"bufio"
	"os"
	"strings"
)

// Ignorer holds a set of key patterns to ignore during diff.
type Ignorer struct {
	patterns []string
}

// Load reads an .envdiffignore file from the given path and returns
// an Ignorer. If the file does not exist, an empty Ignorer is returned.
func Load(path string) (*Ignorer, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Ignorer{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, strings.ToUpper(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Ignorer{patterns: patterns}, nil
}

// Match reports whether the given key should be ignored.
func (ig *Ignorer) Match(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range ig.patterns {
		if upper == p {
			return true
		}
	}
	return false
}

// Filter removes any keys from the provided slice that match ignore patterns.
func (ig *Ignorer) Filter(keys []string) []string {
	if len(ig.patterns) == 0 {
		return keys
	}
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if !ig.Match(k) {
			out = append(out, k)
		}
	}
	return out
}
