package renamer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/renamer"
)

func sampleSuggestions() []renamer.Suggestion {
	return []renamer.Suggestion{
		{From: "DB_HOST", To: "DB_HOSTNAME", Score: 0.80},
		{From: "APP_KEY", To: "APP_SECRET_KEY", Score: 0.67},
	}
}

func TestWriteReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := renamer.WriteReport(&buf, sampleSuggestions(), renamer.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
	if !strings.Contains(out, "DB_HOSTNAME") {
		t.Error("expected DB_HOSTNAME in text output")
	}
	if !strings.Contains(out, "0.80") {
		t.Error("expected score 0.80 in text output")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := renamer.WriteReport(&buf, sampleSuggestions(), renamer.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got []struct {
		From  string  `json:"from"`
		To    string  `json:"to"`
		Score float64 `json:"score"`
	}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].From != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", got[0].From)
	}
}

func TestWriteReport_EmptySuggestions_Text(t *testing.T) {
	var buf bytes.Buffer
	err := renamer.WriteReport(&buf, nil, renamer.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No rename suggestions") {
		t.Error("expected empty message for no suggestions")
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	err := renamer.WriteReport(&buf, sampleSuggestions(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.HasPrefix(strings.TrimSpace(buf.String()), "[") {
		t.Error("default format should not produce JSON")
	}
}
