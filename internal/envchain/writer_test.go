package envchain_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envchain"
)

var sampleResults = []envchain.Result{
	{
		Key:    "DB_URL",
		Value:  "prod-db",
		Source: "prod.env",
		Overridden: []envchain.Override{
			{Source: "base.env", Value: "base-db"},
		},
	},
	{Key: "PORT", Value: "8080", Source: "base.env"},
}

func TestWriteReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := envchain.WriteReport(&buf, sampleResults, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
	if !strings.Contains(out, "prod-db") {
		t.Error("expected prod-db in output")
	}
	if !strings.Contains(out, "overrides") {
		t.Error("expected override line in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in output")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := envchain.WriteReport(&buf, sampleResults, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Errorf("expected 2 entries, got %d", len(decoded))
	}
	if decoded[0]["key"] != "DB_URL" {
		t.Errorf("expected first key DB_URL, got %v", decoded[0]["key"])
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	if err := envchain.WriteReport(&buf, sampleResults, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.HasPrefix(buf.String(), "[") {
		t.Error("default format should not be JSON")
	}
}

func TestWriteReport_EmptyResults_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := envchain.WriteReport(&buf, nil, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no keys resolved") {
		t.Error("expected empty message")
	}
}
