package snapshotter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/snapshotter"
)

func TestWriteReport_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1"}
	before := snapshotter.Capture("staging", env)
	after := snapshotter.Capture("prod", env)
	deltas := snapshotter.Diff(before, after)

	var sb strings.Builder
	snapshotter.WriteReport(&sb, before, after, deltas)
	out := sb.String()

	if !strings.Contains(out, "No changes detected") {
		t.Errorf("expected no-change message, got:\n%s", out)
	}
	if !strings.Contains(out, "staging") || !strings.Contains(out, "prod") {
		t.Errorf("expected snapshot names in output, got:\n%s", out)
	}
}

func TestWriteReport_ShowsAdded(t *testing.T) {
	before := snapshotter.Capture("a", map[string]string{})
	after := snapshotter.Capture("b", map[string]string{"NEW_KEY": "hello"})
	deltas := snapshotter.Diff(before, after)

	var sb strings.Builder
	snapshotter.WriteReport(&sb, before, after, deltas)
	out := sb.String()

	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}

func TestWriteReport_ShowsRemoved(t *testing.T) {
	before := snapshotter.Capture("a", map[string]string{"OLD": "gone"})
	after := snapshotter.Capture("b", map[string]string{})
	deltas := snapshotter.Diff(before, after)

	var sb strings.Builder
	snapshotter.WriteReport(&sb, before, after, deltas)
	out := sb.String()

	if !strings.Contains(out, "- OLD") {
		t.Errorf("expected removed key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}

func TestWriteReport_ShowsChanged(t *testing.T) {
	before := snapshotter.Capture("a", map[string]string{"HOST": "localhost"})
	after := snapshotter.Capture("b", map[string]string{"HOST": "prod.example.com"})
	deltas := snapshotter.Diff(before, after)

	var sb strings.Builder
	snapshotter.WriteReport(&sb, before, after, deltas)
	out := sb.String()

	if !strings.Contains(out, "~ HOST") {
		t.Errorf("expected changed key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "prod.example.com") {
		t.Errorf("expected old and new values, got:\n%s", out)
	}
}

func TestWriteReport_SortedOutput(t *testing.T) {
	before := snapshotter.Capture("a", map[string]string{"Z": "1", "A": "1"})
	after := snapshotter.Capture("b", map[string]string{"Z": "2", "A": "2"})
	deltas := snapshotter.Diff(before, after)

	var sb strings.Builder
	snapshotter.WriteReport(&sb, before, after, deltas)
	out := sb.String()

	idxA := strings.Index(out, "A")
	idxZ := strings.Index(out, "Z")
	if idxA > idxZ {
		t.Errorf("expected A before Z in sorted output")
	}
}
