package snapshotter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes a human-readable diff report of the given deltas to w.
func WriteReport(w io.Writer, before, after Snapshot, deltas []Delta) {
	fmt.Fprintf(w, "Snapshot diff: %s → %s\n", before.Name, after.Name)
	fmt.Fprintf(w, "Captured: %s → %s\n",
		before.CapturedAt.Format("2006-01-02 15:04:05"),
		after.CapturedAt.Format("2006-01-02 15:04:05"),
	)
	fmt.Fprintln(w, strings.Repeat("-", 48))

	if len(deltas) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	sorted := make([]Delta, len(deltas))
	copy(sorted, deltas)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Kind < sorted[j].Kind
	})

	for _, d := range sorted {
		switch d.Kind {
		case Added:
			fmt.Fprintf(w, "  + %-30s = %s\n", d.Key, d.NewValue)
		case Removed:
			fmt.Fprintf(w, "  - %-30s   (was: %s)\n", d.Key, d.OldValue)
		case Changed:
			fmt.Fprintf(w, "  ~ %-30s   %s → %s\n", d.Key, d.OldValue, d.NewValue)
		}
	}

	added := countKind(sorted, Added)
	removed := countKind(sorted, Removed)
	changed := countKind(sorted, Changed)
	fmt.Fprintln(w, strings.Repeat("-", 48))
	fmt.Fprintf(w, "Summary: +%d added, -%d removed, ~%d changed\n", added, removed, changed)
}

func countKind(deltas []Delta, kind DeltaKind) int {
	n := 0
	for _, d := range deltas {
		if d.Kind == kind {
			n++
		}
	}
	return n
}
