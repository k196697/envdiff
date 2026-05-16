package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/validator"
)

func makeResults(kv map[string]map[string]string) []diff.Result {
	var out []diff.Result
	for key, vals := range kv {
		out = append(out, diff.Result{Key: key, Values: vals})
	}
	return out
}

func TestValidate_NoRules(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"PORT": {"prod": "8080"},
	})
	if v := validator.Validate(results, nil); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"SECRET": {"prod": ""},
	})
	rules := []validator.Rule{{Key: "SECRET", Required: true}}
	v := validator.Validate(results, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Rule != "required" {
		t.Errorf("expected rule=required, got %s", v[0].Rule)
	}
}

func TestValidate_InvalidBool(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"ENABLE_CACHE": {"dev": "yes_please"},
	})
	rules := []validator.Rule{{Key: "ENABLE_CACHE", Kind: "bool"}}
	v := validator.Validate(results, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Rule != "bool" {
		t.Errorf("expected rule=bool, got %s", v[0].Rule)
	}
}

func TestValidate_ValidInt(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"PORT": {"prod": "8080"},
	})
	rules := []validator.Rule{{Key: "PORT", Kind: "int"}}
	if v := validator.Validate(results, rules); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestValidate_InvalidURL(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"API_URL": {"staging": "not-a-url"},
	})
	rules := []validator.Rule{{Key: "API_URL", Kind: "url"}}
	v := validator.Validate(results, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_CaseInsensitiveKey(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"db_host": {"prod": ""},
	})
	rules := []validator.Rule{{Key: "DB_HOST", Required: true}}
	v := validator.Validate(results, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for case-insensitive match, got %d", len(v))
	}
}

func TestValidate_EmptyValueSkipsKindCheck(t *testing.T) {
	results := makeResults(map[string]map[string]string{
		"PORT": {"dev": ""},
	})
	rules := []validator.Rule{{Key: "PORT", Kind: "int"}}
	if v := validator.Validate(results, rules); len(v) != 0 {
		t.Fatalf("empty value should skip kind check, got %d violations", len(v))
	}
}
