// Package grouper organises a slice of diff.Result values into named groups
// based on key prefixes.
//
// Keys are split on a configurable separator (default "_") and the leading
// N segments (default 1) are used as the group prefix.  Keys that have no
// separator are collected under the special "(other)" group.
//
// Example
//
//	groups := grouper.ByPrefix(results, grouper.Options{Depth: 1})
//	for _, g := range groups {
//		fmt.Printf("[%s] %d keys\n", g.Prefix, len(g.Results))
//	}
package grouper
