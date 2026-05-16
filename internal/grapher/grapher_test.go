package grapher

import (
	"bytes"
	"strings"
	"testing"
)

func baseEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		".env": {
			"APP_URL":  "https://$HOST:$PORT",
			"HOST":     "localhost",
			"PORT":     "8080",
			"DB_URL":   "postgres://${DB_USER}:${DB_PASS}@localhost/db",
			"DB_USER":  "admin",
			"DB_PASS":  "secret",
		},
	}
}

func TestBuild_DetectsEdges(t *testing.T) {
	g := Build(baseEnvs())
	if len(g.Edges) == 0 {
		t.Fatal("expected edges, got none")
	}
	hasEdge := func(from, to string) bool {
		for _, e := range g.Edges {
			if e.From == from && e.To == to {
				return true
			}
		}
		return false
	}
	if !hasEdge("APP_URL", "HOST") {
		t.Error("expected edge APP_URL -> HOST")
	}
	if !hasEdge("APP_URL", "PORT") {
		t.Error("expected edge APP_URL -> PORT")
	}
	if !hasEdge("DB_URL", "DB_USER") {
		t.Error("expected edge DB_URL -> DB_USER")
	}
	if !hasEdge("DB_URL", "DB_PASS") {
		t.Error("expected edge DB_URL -> DB_PASS")
	}
}

func TestBuild_NoEdges(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"FOO": "bar", "BAZ": "qux"},
	}
	g := Build(envs)
	if len(g.Edges) != 0 {
		t.Errorf("expected no edges, got %d", len(g.Edges))
	}
}

func TestBuild_OrphansDetected(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {
			"APP_URL": "https://$UNDEFINED_HOST",
		},
	}
	g := Build(envs)
	if len(g.Orphans) != 1 || g.Orphans[0] != "UNDEFINED_HOST" {
		t.Errorf("expected orphan UNDEFINED_HOST, got %v", g.Orphans)
	}
}

func TestBuild_AllKeysSorted(t *testing.T) {
	g := Build(baseEnvs())
	for i := 1; i < len(g.AllKeys); i++ {
		if g.AllKeys[i] < g.AllKeys[i-1] {
			t.Errorf("AllKeys not sorted at index %d: %v", i, g.AllKeys)
		}
	}
}

func TestWrite_NoEdges(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {"FOO": "bar"},
	}
	g := Build(envs)
	var buf bytes.Buffer
	Write(g, &buf)
	if !strings.Contains(buf.String(), "No key references") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWrite_ShowsEdgesAndOrphans(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {
			"URL": "http://$HOST:$MISSING",
			"HOST": "localhost",
		},
	}
	g := Build(envs)
	var buf bytes.Buffer
	Write(g, &buf)
	out := buf.String()
	if !strings.Contains(out, "URL -> HOST") {
		t.Errorf("expected edge in output, got:\n%s", out)
	}
	if !strings.Contains(out, "! MISSING") {
		t.Errorf("expected orphan MISSING in output, got:\n%s", out)
	}
}

func TestExtractRefs_BraceAndBare(t *testing.T) {
	refs := extractRefs("${FOO}_bar_$BAZ")
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %v", refs)
	}
	if refs[0] != "FOO" || refs[1] != "BAZ" {
		t.Errorf("unexpected refs: %v", refs)
	}
}
