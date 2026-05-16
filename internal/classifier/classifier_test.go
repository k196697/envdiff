package classifier_test

import (
	"testing"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/diff"
)

func makeResult(key, status string) diff.Result {
	return diff.Result{Key: key, Status: diff.Status(status)}
}

func TestClassify_AllMatch(t *testing.T) {
	input := []diff.Result{
		makeResult("APP_NAME", "match"),
		makeResult("PORT", "match"),
	}
	got := classifier.Classify(input, classifier.Options{})
	for _, r := range got {
		if r.Severity != classifier.SeverityInfo {
			t.Errorf("key %s: expected info, got %s", r.Diff.Key, r.Severity)
		}
	}
}

func TestClassify_MissingIsWarning(t *testing.T) {
	input := []diff.Result{makeResult("MISSING_KEY", "missing")}
	got := classifier.Classify(input, classifier.Options{})
	if got[0].Severity != classifier.SeverityWarning {
		t.Errorf("expected warning, got %s", got[0].Severity)
	}
}

func TestClassify_CriticalPrefix(t *testing.T) {
	input := []diff.Result{makeResult("SECRET_TOKEN", "missing")}
	opts := classifier.Options{CriticalPrefixes: []string{"SECRET_"}}
	got := classifier.Classify(input, opts)
	if got[0].Severity != classifier.SeverityCritical {
		t.Errorf("expected critical, got %s", got[0].Severity)
	}
}

func TestClassify_CriticalPrefixCaseInsensitive(t *testing.T) {
	input := []diff.Result{makeResult("secret_token", "match")}
	opts := classifier.Options{CriticalPrefixes: []string{"SECRET_"}}
	got := classifier.Classify(input, opts)
	if got[0].Severity != classifier.SeverityCritical {
		t.Errorf("expected critical, got %s", got[0].Severity)
	}
}

func TestClassify_WarningPrefixOverridesDefault(t *testing.T) {
	input := []diff.Result{makeResult("DB_HOST", "match")}
	opts := classifier.Options{WarningPrefixes: []string{"DB_"}}
	got := classifier.Classify(input, opts)
	if got[0].Severity != classifier.SeverityWarning {
		t.Errorf("expected warning, got %s", got[0].Severity)
	}
}

func TestClassify_CriticalTakesPriorityOverWarning(t *testing.T) {
	input := []diff.Result{makeResult("SECRET_DB_PASS", "mismatch")}
	opts := classifier.Options{
		CriticalPrefixes: []string{"SECRET_"},
		WarningPrefixes:  []string{"SECRET_DB_"},
	}
	got := classifier.Classify(input, opts)
	if got[0].Severity != classifier.SeverityCritical {
		t.Errorf("expected critical, got %s", got[0].Severity)
	}
}

func TestClassify_EmptyInput(t *testing.T) {
	got := classifier.Classify(nil, classifier.Options{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}
