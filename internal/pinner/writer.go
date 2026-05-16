package pinner

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteReport writes the list of currently pinned keys to w in the requested
// format ("text" or "json"). An empty format defaults to "text".
func (ps *PinnedSet) WriteReport(w io.Writer, format string) error {
	keys := make([]string, 0, len(ps.keys))
	for k := range ps.keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case "json":
		return writeJSON(w, keys)
	default:
		return writeText(w, keys)
	}
}

func writeText(w io.Writer, keys []string) error {
	if len(keys) == 0 {
		_, err := fmt.Fprintln(w, "No pinned keys.")
		return err
	}
	_, err := fmt.Fprintln(w, "Pinned keys (mismatches suppressed):")
	if err != nil {
		return err
	}
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "  - %s\n", k); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, keys []string) error {
	payload := map[string]interface{}{
		"pinned_keys": keys,
		"count":       len(keys),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
