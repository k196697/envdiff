// Package pinner allows users to pin specific keys so their values are
// treated as authoritative and never flagged as mismatches during comparison.
package pinner

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// PinnedSet holds the set of keys whose values are considered authoritative.
type PinnedSet struct {
	keys map[string]struct{}
}

// Load reads a pinned-keys file (one key per line, # comments allowed).
// If path is empty or the file does not exist, an empty PinnedSet is returned.
func Load(path string) (*PinnedSet, error) {
	ps := &PinnedSet{keys: make(map[string]struct{})}
	if path == "" {
		return ps, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ps, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ps.keys[strings.ToUpper(line)] = struct{}{}
	}
	return ps, scanner.Err()
}

// IsPinned reports whether key is in the pinned set.
func (ps *PinnedSet) IsPinned(key string) bool {
	_, ok := ps.keys[strings.ToUpper(key)]
	return ok
}

// Apply filters results, removing mismatch entries for pinned keys.
// Missing entries are preserved — pinning only suppresses value mismatches.
func (ps *PinnedSet) Apply(results []diff.Result) []diff.Result {
	out := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if r.Status == diff.Mismatch && ps.IsPinned(r.Key) {
			continue
		}
		out = append(out, r)
	}
	return out
}
