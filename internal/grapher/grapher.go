// Package grapher builds a dependency graph showing which keys reference
// other keys across multiple .env files.
package grapher

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Edge represents a directed dependency from one key to another.
type Edge struct {
	From string // key that contains the reference
	To   string // key being referenced
	File string // source file where From was found
}

// Graph holds all detected edges and the full set of known keys.
type Graph struct {
	Edges    []Edge
	AllKeys  []string
	Orphans  []string // keys referenced but never defined
}

// Build constructs a Graph from a map of filename → key/value pairs.
// It detects references of the form $KEY or ${KEY} in values.
func Build(envs map[string]map[string]string) Graph {
	known := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			known[k] = struct{}{}
		}
	}

	var edges []Edge
	referenced := make(map[string]struct{})

	for file, kv := range envs {
		for key, val := range kv {
			for _, ref := range extractRefs(val) {
				referenced[ref] = struct{}{}
				edges = append(edges, Edge{From: key, To: ref, File: file})
			}
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].File != edges[j].File {
			return edges[i].File < edges[j].File
		}
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		return edges[i].To < edges[j].To
	})

	allKeys := make([]string, 0, len(known))
	for k := range known {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	var orphans []string
	for ref := range referenced {
		if _, ok := known[ref]; !ok {
			orphanss = append(orphans, ref)
		}
	}
	sort.Strings(orphans)

	return Graph{Edges: edges, AllKeys: allKeys, Orphans: orphans}
}

// Write renders the graph as a human-readable report to w.
func Write(g Graph, w io.Writer) {
	if len(g.Edges) == 0 {
		fmt.Fprintln(w, "No key references detected.")
		return
	}
	fmt.Fprintf(w, "Key dependency graph (%d edge(s)):\n", len(g.Edges))
	for _, e := range g.Edges {
		fmt.Fprintf(w, "  [%s] %s -> %s\n", e.File, e.From, e.To)
	}
	if len(g.Orphans) > 0 {
		fmt.Fprintln(w, "\nUndefined references:")
		for _, o := range g.Orphans {
			fmt.Fprintf(w, "  ! %s\n", o)
		}
	}
}

// extractRefs returns all key names referenced via $KEY or ${KEY} in s.
func extractRefs(s string) []string {
	var refs []string
	for i := 0; i < len(s); i++ {
		if s[i] != '$' {
			continue
		}
		rest := s[i+1:]
		if strings.HasPrefix(rest, "{") {
			end := strings.Index(rest, "}")
			if end > 1 {
				refs = append(refs, rest[1:end])
				i += end + 1
			}
		} else {
			j := 0
			for j < len(rest) && isKeyChar(rest[j]) {
				j++
			}
			if j > 0 {
				refs = append(refs, rest[:j])
				i += j
			}
		}
	}
	return refs
}

func isKeyChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}
