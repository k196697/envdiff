package pinner_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/pinner"
)

func loadPinnedSet(t *testing.T, keys ...string) *pinner.PinnedSet {
	t.Helper()
	content := strings.Join(keys, "\n") + "\n"
	p := writeTempPins(t, content)
	ps, err := pinner.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	return ps
}

func TestWriteReport_TextFormat(t *testing.T) {
	ps := loadPinnedSet(t, "DB_HOST", "API_KEY")
	var buf bytes.Buffer
	if err := ps.WriteReport(&buf, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
	if !strings.Contains(out, "Pinned keys") {
		t.Error("expected header in text output")
	}
}

func TestWriteReport_TextFormat_Empty(t *testing.T) {
	ps, _ := pinner.Load("")
	var buf bytes.Buffer
	if err := ps.WriteReport(&buf, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No pinned keys") {
		t.Error("expected empty message")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	ps := loadPinnedSet(t, "PORT", "DB_HOST")
	var buf bytes.Buffer
	if err := ps.WriteReport(&buf, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	count, ok := payload["count"].(float64)
	if !ok || int(count) != 2 {
		t.Errorf("expected count 2, got %v", payload["count"])
	}
	keys, ok := payload["pinned_keys"].([]interface{})
	if !ok || len(keys) != 2 {
		t.Errorf("expected 2 pinned_keys, got %v", payload["pinned_keys"])
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	ps := loadPinnedSet(t, "FOO")
	var buf bytes.Buffer
	if err := ps.WriteReport(&buf, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "FOO") {
		t.Error("expected FOO in default (text) output")
	}
}
