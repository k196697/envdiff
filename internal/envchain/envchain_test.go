package envchain_test

import (
	"testing"

	"github.com/user/envdiff/internal/envchain"
)

func makeLink(name string, pairs ...string) envchain.Link {
	env := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return envchain.Link{Name: name, Env: env}
}

func TestResolve_SingleLink(t *testing.T) {
	chain := []envchain.Link{makeLink("base", "FOO", "bar", "BAZ", "qux")}
	results := envchain.Resolve(chain)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestResolve_FirstLinkWins(t *testing.T) {
	chain := []envchain.Link{
		makeLink("prod", "DB_URL", "prod-db"),
		makeLink("base", "DB_URL", "base-db"),
	}
	results := envchain.Resolve(chain)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Value != "prod-db" {
		t.Errorf("expected prod-db, got %s", r.Value)
	}
	if r.Source != "prod" {
		t.Errorf("expected source prod, got %s", r.Source)
	}
	if len(r.Overridden) != 1 || r.Overridden[0].Value != "base-db" {
		t.Errorf("expected one override with base-db")
	}
}

func TestResolve_MergesUniqueKeys(t *testing.T) {
	chain := []envchain.Link{
		makeLink("a", "ALPHA", "1"),
		makeLink("b", "BETA", "2"),
		makeLink("c", "GAMMA", "3"),
	}
	results := envchain.Resolve(chain)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestResolve_SortedByKey(t *testing.T) {
	chain := []envchain.Link{
		makeLink("x", "ZEBRA", "z", "ALPHA", "a", "MANGO", "m"),
	}
	results := envchain.Resolve(chain)
	keys := []string{results[0].Key, results[1].Key, results[2].Key}
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("results not sorted: %v", keys)
	}
}

func TestFlatten(t *testing.T) {
	results := []envchain.Result{
		{Key: "A", Value: "1", Source: "f1"},
		{Key: "B", Value: "2", Source: "f2"},
	}
	m := envchain.Flatten(results)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("flatten mismatch: %v", m)
	}
}

func TestValidate_Empty(t *testing.T) {
	if err := envchain.Validate(nil); err == nil {
		t.Error("expected error for empty chain")
	}
}

func TestValidate_BlankName(t *testing.T) {
	chain := []envchain.Link{{Name: "", Env: map[string]string{"X": "1"}}}
	if err := envchain.Validate(chain); err == nil {
		t.Error("expected error for blank link name")
	}
}

func TestValidate_Valid(t *testing.T) {
	chain := []envchain.Link{makeLink("base", "K", "v")}
	if err := envchain.Validate(chain); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
