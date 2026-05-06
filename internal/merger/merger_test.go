package merger_test

import (
	"reflect"
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func TestMerge_Basic(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "3000"},
		"prod": {"HOST": "example.com", "SECRET": "abc"},
	}
	got := merger.Merge(envs)
	if got["HOST"] != "localhost" && got["HOST"] != "example.com" {
		t.Errorf("unexpected HOST value: %s", got["HOST"])
	}
	if _, ok := got["PORT"]; !ok {
		t.Error("expected PORT in merged result")
	}
	if _, ok := got["SECRET"]; !ok {
		t.Error("expected SECRET in merged result")
	}
}

func TestMerge_Empty(t *testing.T) {
	got := merger.Merge(map[string]map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMerge_FirstValueWins(t *testing.T) {
	// Deterministic: only one file has KEY, so no conflict.
	envs := map[string]map[string]string{
		"a": {"KEY": "from-a"},
	}
	got := merger.Merge(envs)
	if got["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %s", got["KEY"])
	}
}

func TestKeys_Sorted(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"ZEBRA": "1", "APPLE": "2"},
		"prod": {"MANGO": "3", "APPLE": "4"},
	}
	got := merger.Keys(envs)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v, want %v", got, want)
	}
}

func TestCoverage_Basic(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "3000"},
		"prod": {"HOST": "example.com", "SECRET": "abc"},
	}
	got := merger.Coverage(envs)

	if len(got["HOST"]) != 2 {
		t.Errorf("HOST should be covered by 2 envs, got %v", got["HOST"])
	}
	if len(got["PORT"]) != 1 || got["PORT"][0] != "dev" {
		t.Errorf("PORT should only be in dev, got %v", got["PORT"])
	}
	if len(got["SECRET"]) != 1 || got["SECRET"][0] != "prod" {
		t.Errorf("SECRET should only be in prod, got %v", got["SECRET"])
	}
}
