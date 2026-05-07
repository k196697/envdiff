package linter

import (
	"strings"
	"testing"
)

func TestLint_CleanEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"API_KEY":      "abc123",
	}
	issues := Lint("prod.env", env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "",
	}
	issues := Lint("dev.env", env)
	if !hasMsg(issues, "value is empty") {
		t.Error("expected 'value is empty' issue")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
	}
	issues := Lint("dev.env", env)
	if !hasMsg(issues, "key is not uppercase") {
		t.Error("expected 'key is not uppercase' issue")
	}
}

func TestLint_KeyWithSpace(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
	}
	issues := Lint("dev.env", env)
	if !hasMsg(issues, "key contains whitespace") {
		t.Error("expected 'key contains whitespace' issue")
	}
}

func TestLint_UnresolvedPlaceholder(t *testing.T) {
	cases := []string{"${SOME_VAR}", "<YOUR_SECRET>", "changeme", "TODO", "FIXME"}
	for _, val := range cases {
		env := map[string]string{"SOME_KEY": val}
		issues := Lint("test.env", env)
		if !hasMsg(issues, "unresolved placeholder") {
			t.Errorf("expected unresolved placeholder issue for value %q", val)
		}
	}
}

func TestIssue_String(t *testing.T) {
	i := Issue{File: "prod.env", Key: "FOO", Msg: "value is empty"}
	got := i.String()
	if !strings.Contains(got, "prod.env") || !strings.Contains(got, "FOO") {
		t.Errorf("unexpected String() output: %s", got)
	}
}

// hasMsg returns true if any issue message contains the given substring.
func hasMsg(issues []Issue, substr string) bool {
	for _, i := range issues {
		if strings.Contains(i.Msg, substr) {
			return true
		}
	}
	return false
}
