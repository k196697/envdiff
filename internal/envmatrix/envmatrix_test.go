package envmatrix_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/envmatrix"
)

var sampleEnvs = map[string]map[string]string{
	"production": {"DB_HOST": "prod-db", "API_KEY": "abc123", "LOG_LEVEL": "error"},
	"staging":    {"DB_HOST": "stage-db", "API_KEY": "xyz789"},
	"development": {"DB_HOST": "localhost", "DEBUG": "true", "LOG_LEVEL": "debug"},
}

func TestBuild_EnvsAreSorted(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	if len(m.Envs) != 3 {
		t.Fatalf("expected 3 envs, got %d", len(m.Envs))
	}
	if m.Envs[0] != "development" || m.Envs[1] != "production" || m.Envs[2] != "staging" {
		t.Errorf("envs not sorted: %v", m.Envs)
	}
}

func TestBuild_KeysAreSorted(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	keys := make([]string, len(m.Rows))
	for i, r := range m.Rows {
		keys[i] = r.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted at index %d: %v", i, keys)
		}
	}
}

func TestBuild_AbsentTracked(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	var debugRow *envmatrix.Row
	for i := range m.Rows {
		if m.Rows[i].Key == "DEBUG" {
			debugRow = &m.Rows[i]
			break
		}
	}
	if debugRow == nil {
		t.Fatal("DEBUG key not found in matrix")
	}
	if len(debugRow.Absent) != 2 {
		t.Errorf("expected DEBUG absent in 2 envs, got %d", len(debugRow.Absent))
	}
	if _, ok := debugRow.Values["development"]; !ok {
		t.Error("expected DEBUG present in development")
	}
}

func TestBuild_EmptyEnvs(t *testing.T) {
	m := envmatrix.Build(map[string]map[string]string{})
	if len(m.Rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(m.Rows))
	}
}

func TestWrite_TextFormat(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	var buf bytes.Buffer
	if err := envmatrix.Write(m, "text", &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
	if !strings.Contains(out, "<missing>") {
		t.Error("expected <missing> marker in text output")
	}
	if !strings.Contains(out, "production") {
		t.Error("expected env name in header")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	var buf bytes.Buffer
	if err := envmatrix.Write(m, "json", &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var decoded envmatrix.Matrix
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}
	if len(decoded.Rows) != len(m.Rows) {
		t.Errorf("expected %d rows, got %d", len(m.Rows), len(decoded.Rows))
	}
}

func TestWrite_DefaultIsText(t *testing.T) {
	m := envmatrix.Build(sampleEnvs)
	var buf bytes.Buffer
	if err := envmatrix.Write(m, "", &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if strings.HasPrefix(strings.TrimSpace(buf.String()), "{") {
		t.Error("default format should be text, not JSON")
	}
}

func TestWrite_EmptyMatrix_Text(t *testing.T) {
	m := envmatrix.Build(map[string]map[string]string{})
	var buf bytes.Buffer
	if err := envmatrix.Write(m, "text", &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "no keys found") {
		t.Error("expected 'no keys found' for empty matrix")
	}
}
