// Package merger provides functionality to merge multiple parsed env maps
// into a unified superset, useful for generating a canonical .env.example.
package merger

import "sort"

// Merge combines multiple named env maps into a single superset map.
// Keys present in any of the input maps will appear in the result.
// If the same key appears in multiple maps, the value from the first
// map (in input order) that defines it is used.
func Merge(envs map[string]map[string]string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}

// Keys returns a sorted slice of all unique keys across the provided env maps.
func Keys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Coverage returns, for each key, the set of environment names that define it.
func Coverage(envs map[string]map[string]string) map[string][]string {
	result := make(map[string][]string)
	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		env := envs[name]
		for k := range env {
			result[k] = append(result[k], name)
		}
	}
	return result
}
