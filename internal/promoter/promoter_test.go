package promoter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/promoter"
)

var stagingEnv = map[string]string{
	"DB_HOST":     "staging-db",
	"DB_PORT":     "5432",
	"FEATURE_X":  "true",
	"LOG_LEVEL":  "debug",
}

var productionEnv = map[string]string{
	"DB_HOST": "prod-db",
	"DB_PORT": "5432",
}

func TestPromote_OnlyMissing(t *testing.T) {
	opts := promoter.DefaultOptions()
	suggestions := promoter.Promote(stagingEnv, productionEnv, "staging", "production", opts)

	if len(suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
	}
	keys := make([]string, len(suggestions))
	for i, s := range suggestions {
		keys[i] = s.Key
	}
	if keys[0] != "FEATURE_X" || keys[1] != "LOG_LEVEL" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestPromote_IncludeMismatch(t *testing.T) {
	opts := promoter.Options{OnlyMissing: false, IncludeMismatch: true}
	suggestions := promoter.Promote(stagingEnv, productionEnv, "staging", "production", opts)

	// DB_HOST differs; FEATURE_X and LOG_LEVEL are missing but OnlyMissing=false
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion (mismatch only), got %d", len(suggestions))
	}
	if suggestions[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", suggestions[0].Key)
	}
}

func TestPromote_NoSuggestions(t *testing.T) {
	same := map[string]string{"A": "1"}
	opts := promoter.DefaultOptions()
	suggestions := promoter.Promote(same, same, "env1", "env2", opts)
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions, got %d", len(suggestions))
	}
}

func TestWriteReport_TextFormat(t *testing.T) {
	suggestions := []promoter.Suggestion{
		{Key: "LOG_LEVEL", FromEnv: "staging", FromValue: "debug", ToEnv: "production", Reason: "key missing in production"},
	}
	var buf bytes.Buffer
	if err := promoter.WriteReport(&buf, suggestions, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL in output, got: %s", out)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in output")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	suggestions := []promoter.Suggestion{
		{Key: "FEATURE_X", FromEnv: "staging", FromValue: "true", ToEnv: "production", Reason: "key missing in production"},
	}
	var buf bytes.Buffer
	if err := promoter.WriteReport(&buf, suggestions, "json"); err != nil {
		t.Fatal(err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out[0]["key"] != "FEATURE_X" {
		t.Errorf("expected FEATURE_X, got %s", out[0]["key"])
	}
}

func TestWriteReport_EmptySuggestions(t *testing.T) {
	var buf bytes.Buffer
	if err := promoter.WriteReport(&buf, nil, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No promotion candidates") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
