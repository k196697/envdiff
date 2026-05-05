package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/filter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestApply_NoOptions(t *testing.T) {
	env := baseEnv()
	result := filter.Apply(env, filter.Options{})
	if len(result) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{Prefix: "APP_"})
	if len(result) != 2 {
		t.Errorf("expected 2 keys with prefix APP_, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestApply_PrefixCaseInsensitive(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{Prefix: "app_"})
	if len(result) != 2 {
		t.Errorf("expected 2 keys with lowercase prefix app_, got %d", len(result))
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{
		Exclude: []string{"DB_PASSWORD", "LOG_LEVEL"},
	})
	if len(result) != 3 {
		t.Errorf("expected 3 keys after exclusion, got %d", len(result))
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
	if _, ok := result["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should have been excluded")
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{
		Prefix:  "DB_",
		Exclude: []string{"DB_PASSWORD"},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	result := filter.Apply(map[string]string{}, filter.Options{Prefix: "APP_"})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
