// Package envchain resolves a chain of .env files in priority order,
// producing a final merged map where earlier files take precedence.
package envchain

import (
	"fmt"
	"sort"
)

// Link represents a single file in the chain with its parsed key/value pairs.
type Link struct {
	Name string
	Env  map[string]string
}

// Result holds the resolved value for a key across the chain.
type Result struct {
	Key      string
	Value    string
	Source   string // name of the file that provided the value
	Overridden []Override
}

// Override records a value that was shadowed by a higher-priority link.
type Override struct {
	Source string
	Value  string
}

// Resolve walks the chain from highest to lowest priority and returns
// one Result per unique key. The first link that defines a key wins.
func Resolve(chain []Link) []Result {
	seen := make(map[string]*Result)
	var order []string

	for _, link := range chain {
		for k, v := range link.Env {
			if r, exists := seen[k]; !exists {
				seen[k] = &Result{
					Key:    k,
					Value:  v,
					Source: link.Name,
				}
				order = append(order, k)
			} else {
				r.Overridden = append(r.Overridden, Override{
					Source: link.Name,
					Value:  v,
				})
			}
		}
	}

	sort.Strings(order)
	results := make([]Result, 0, len(order))
	for _, k := range order {
		results = append(results, *seen[k])
	}
	return results
}

// Flatten converts Resolve output into a plain key→value map.
func Flatten(results []Result) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		out[r.Key] = r.Value
	}
	return out
}

// Validate checks that the chain contains at least one link and that no
// link has a blank name, returning a descriptive error when violated.
func Validate(chain []Link) error {
	if len(chain) == 0 {
		return fmt.Errorf("envchain: chain must contain at least one link")
	}
	for i, l := range chain {
		if l.Name == "" {
			return fmt.Errorf("envchain: link at index %d has an empty name", i)
		}
	}
	return nil
}
