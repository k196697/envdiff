package envnormaliser

import (
	"testing"
)

func TestApply_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "App_Port": "8080"}
	got := Apply(env, Options{UppercaseKeys: true})
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
	if _, ok := got["APP_PORT"]; !ok {
		t.Error("expected APP_PORT key")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestApply_NoUppercase(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	got := Apply(env, Options{UppercaseKeys: false})
	if _, ok := got["db_host"]; !ok {
		t.Error("expected original key preserved")
	}
}

func TestApply_TrimValues(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	got := Apply(env, Options{TrimValues: true})
	if got["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", got["KEY"])
	}
}

func TestApply_CollapseValues(t *testing.T) {
	env := map[string]string{"MSG": "hello   world  foo"}
	got := Apply(env, Options{CollapseValues: true})
	if got["MSG"] != "hello world foo" {
		t.Errorf("unexpected collapsed value: %q", got["MSG"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"key": "  val  "}
	_ = Apply(env, DefaultOptions())
	if env["key"] != "  val  " {
		t.Error("original map was mutated")
	}
}

func TestApplyAll_NormalisesAllEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		".env.staging": {"db_url": " postgres://localhost "},
		".env.prod":    {"db_url": " postgres://prod "},
	}
	got := ApplyAll(envs, DefaultOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(got))
	}
	for name, env := range got {
		if _, ok := env["DB_URL"]; !ok {
			t.Errorf("%s: expected DB_URL key", name)
		}
		val := env["DB_URL"]
		if len(val) == 0 || val[0] == ' ' || val[len(val)-1] == ' ' {
			t.Errorf("%s: value not trimmed: %q", name, val)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if !opts.UppercaseKeys {
		t.Error("expected UppercaseKeys true by default")
	}
	if !opts.TrimValues {
		t.Error("expected TrimValues true by default")
	}
	if opts.CollapseValues {
		t.Error("expected CollapseValues false by default")
	}
}
