package cascader_test

import (
	"testing"

	"github.com/user/envdiff/internal/cascader"
)

func TestCascade_Empty(t *testing.T) {
	results := cascader.Cascade(nil)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestCascade_SingleLayer(t *testing.T) {
	layers := []cascader.Layer{
		{Name: "prod", Env: map[string]string{"DB_HOST": "prod.db", "PORT": "5432"}},
	}
	results := cascader.Cascade(layers)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.SourceName != "prod" {
			t.Errorf("key %q: expected source prod, got %q", r.Key, r.SourceName)
		}
		if len(r.Overridden) != 0 {
			t.Errorf("key %q: unexpected overrides", r.Key)
		}
	}
}

func TestCascade_HigherPriorityWins(t *testing.T) {
	layers := []cascader.Layer{
		{Name: "local", Env: map[string]string{"DB_HOST": "localhost"}},
		{Name: "prod", Env: map[string]string{"DB_HOST": "prod.db"}},
	}
	results := cascader.Cascade(layers)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Value != "localhost" {
		t.Errorf("expected value localhost, got %q", r.Value)
	}
	if r.SourceName != "local" {
		t.Errorf("expected source local, got %q", r.SourceName)
	}
	if len(r.Overridden) != 1 || r.Overridden[0].LayerName != "prod" {
		t.Errorf("expected one override from prod, got %+v", r.Overridden)
	}
}

func TestCascade_KeyOnlyInLowerLayer(t *testing.T) {
	layers := []cascader.Layer{
		{Name: "local", Env: map[string]string{"APP_ENV": "dev"}},
		{Name: "prod", Env: map[string]string{"APP_ENV": "prod", "SECRET": "xyz"}},
	}
	results := cascader.Cascade(layers)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	byKey := make(map[string]cascader.Result)
	for _, r := range results {
		byKey[r.Key] = r
	}
	if byKey["SECRET"].SourceName != "prod" {
		t.Errorf("SECRET should come from prod, got %q", byKey["SECRET"].SourceName)
	}
	if len(byKey["SECRET"].Overridden) != 0 {
		t.Error("SECRET should have no overrides")
	}
}

func TestCascade_ResultsAreSortedByKey(t *testing.T) {
	layers := []cascader.Layer{
		{Name: "a", Env: map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}},
	}
	results := cascader.Cascade(layers)
	for i := 1; i < len(results); i++ {
		if results[i].Key < results[i-1].Key {
			t.Errorf("results not sorted: %q before %q", results[i-1].Key, results[i].Key)
		}
	}
}
