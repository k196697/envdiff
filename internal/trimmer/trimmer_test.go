package trimmer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/trimmer"
)

func makeResults(vals map[string]map[string]string) []diff.Result {
	var results []diff.Result
	for key, fileVals := range vals {
		results = append(results, diff.Result{
			Key:    key,
			Values: fileVals,
		})
	}
	return results
}

func TestDetect_NoIssues(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"DB_HOST": {"a.env": "localhost", "b.env": "prod-db"},
	})
	issues := trimmer.Detect(results)
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(issues))
	}
}

func TestDetect_TrailingSpace(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"API_KEY": {"a.env": "secret "},
	})
	issues := trimmer.Detect(results)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].Trailing {
		t.Error("expected Trailing=true")
	}
	if issues[0].Leading {
		t.Error("expected Leading=false")
	}
}

func TestDetect_LeadingSpace(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"DB_PASS": {"b.env": " hunter2"},
	})
	issues := trimmer.Detect(results)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].Leading {
		t.Error("expected Leading=true")
	}
}

func TestDetect_BothSides(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"TOKEN": {"c.env": "\t abc \t"},
	})
	issues := trimmer.Detect(results)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !issues[0].Leading || !issues[0].Trailing {
		t.Error("expected both Leading and Trailing true")
	}
}

func TestDetect_SortedOutput(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"Z_KEY": {"a.env": "val "},
		"A_KEY": {"a.env": " val"},
	})
	issues := trimmer.Detect(results)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Key != "A_KEY" || issues[1].Key != "Z_KEY" {
		t.Errorf("unexpected order: %s, %s", issues[0].Key, issues[1].Key)
	}
}

func TestWriteReport_TextNoIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := trimmer.WriteReport(&buf, nil, "text"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no whitespace issues") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteReport_TextWithIssues(t *testing.T) {
	issues := []trimmer.Issue{
		{Key: "FOO", File: "dev.env", Value: "bar ", Trailing: true},
	}
	var buf bytes.Buffer
	if err := trimmer.WriteReport(&buf, issues, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[TRIM]") || !strings.Contains(out, "FOO") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestWriteReport_JSONFormat(t *testing.T) {
	issues := []trimmer.Issue{
		{Key: "BAR", File: "prod.env", Value: " baz", Leading: true},
	}
	var buf bytes.Buffer
	if err := trimmer.WriteReport(&buf, issues, "json"); err != nil {
		t.Fatal(err)
	}
	var out []trimmer.Issue
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 || out[0].Key != "BAR" {
		t.Errorf("unexpected decoded output: %+v", out)
	}
}

func TestWriteReport_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	if err := trimmer.WriteReport(&buf, nil, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no whitespace issues") {
		t.Errorf("default format should be text, got: %s", buf.String())
	}
}
