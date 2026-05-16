package envprinter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envprinter"
)

var sampleEnv = map[string]string{
	"APP_NAME": "envdiff",
	"DEBUG":    "true",
	"SECRET":   "s3cr3t",
	"EMPTY":    "",
}

func TestWrite_TextFormat_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	opts := envprinter.DefaultOptions()
	if err := envprinter.Write(&buf, sampleEnv, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != len(sampleEnv) {
		t.Fatalf("expected %d lines, got %d", len(sampleEnv), len(lines))
	}
	// First key alphabetically should be APP_NAME
	if !strings.HasPrefix(lines[0], "APP_NAME=") {
		t.Errorf("expected first line to start with APP_NAME=, got %q", lines[0])
	}
}

func TestWrite_TextFormat_MaskValues(t *testing.T) {
	var buf bytes.Buffer
	opts := envprinter.DefaultOptions()
	opts.MaskValues = true
	if err := envprinter.Write(&buf, sampleEnv, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if strings.Contains(output, "s3cr3t") {
		t.Error("expected secret value to be masked")
	}
	if !strings.Contains(output, "***") {
		t.Error("expected masked value '***' in output")
	}
	// Empty values should not be masked
	if strings.Contains(output, "EMPTY=***") {
		t.Error("empty values should not be replaced with ***")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := envprinter.Options{Format: "json", SortKeys: true}
	if err := envprinter.Write(&buf, sampleEnv, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["APP_NAME"] != "envdiff" {
		t.Errorf("expected APP_NAME=envdiff, got %q", result["APP_NAME"])
	}
}

func TestWrite_JSONFormat_MaskValues(t *testing.T) {
	var buf bytes.Buffer
	opts := envprinter.Options{Format: "json", MaskValues: true}
	if err := envprinter.Write(&buf, sampleEnv, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["SECRET"] != "***" {
		t.Errorf("expected SECRET to be masked, got %q", result["SECRET"])
	}
	if result["EMPTY"] != "" {
		t.Errorf("expected EMPTY to remain empty, got %q", result["EMPTY"])
	}
}

func TestWrite_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	opts := envprinter.Options{SortKeys: true}
	if err := envprinter.Write(&buf, sampleEnv, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty text output")
	}
	if strings.HasPrefix(buf.String(), "{") {
		t.Error("expected text output, not JSON")
	}
}
