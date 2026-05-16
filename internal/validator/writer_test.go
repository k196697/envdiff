package validator_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/validator"
)

var sampleViolations = []validator.Violation{
	{File: "prod.env", Key: "PORT", Rule: "int", Message: "int: \"abc\" is not an integer"},
	{File: "dev.env", Key: "SECRET", Rule: "required", Message: "value is required but empty"},
}

func TestWriteReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := validator.WriteReport(&buf, sampleViolations, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in text output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET in text output")
	}
}

func TestWriteReport_TextFormat_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	if err := validator.WriteReport(&buf, nil, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Error("expected 'no violations' message")
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := validator.WriteReport(&buf, sampleViolations, "json"); err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Violations []validator.Violation `json:"violations"`
		Total       int                   `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Total != 2 {
		t.Errorf("expected total=2, got %d", payload.Total)
	}
	if len(payload.Violations) != 2 {
		t.Errorf("expected 2 violations in JSON, got %d", len(payload.Violations))
	}
}

func TestWriteReport_JSONFormat_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := validator.WriteReport(&buf, nil, "json"); err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Violations []validator.Violation `json:"violations"`
		Total       int                   `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Total != 0 {
		t.Errorf("expected total=0, got %d", payload.Total)
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	if err := validator.WriteReport(&buf, sampleViolations, ""); err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(buf.String(), "{") {
		t.Error("default format should be text, not JSON")
	}
}
