package resolver_test

import (
	"sort"
	"testing"

	"github.com/user/envdiff/internal/resolver"
)

func resultByKey(results []resolver.Result, key string) (resolver.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return resolver.Result{}, false
}

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	results := resolver.Resolve(env)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	r, ok := resultByKey(results, "FOO")
	if !ok || r.Resolved != "bar" {
		t.Errorf("expected FOO=bar, got %q", r.Resolved)
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "URL": "http://${HOST}:8080"}
	results := resolver.Resolve(env)
	r, ok := resultByKey(results, "URL")
	if !ok {
		t.Fatal("URL result not found")
	}
	if r.Resolved != "http://localhost:8080" {
		t.Errorf("unexpected resolved value: %q", r.Resolved)
	}
	if len(r.Refs) != 1 || r.Refs[0] != "HOST" {
		t.Errorf("expected refs [HOST], got %v", r.Refs)
	}
}

func TestResolve_DollarStyle(t *testing.T) {
	env := map[string]string{"NAME": "world", "GREETING": "hello $NAME"}
	results := resolver.Resolve(env)
	r, ok := resultByKey(results, "GREETING")
	if !ok {
		t.Fatal("GREETING result not found")
	}
	if r.Resolved != "hello world" {
		t.Errorf("unexpected resolved value: %q", r.Resolved)
	}
}

func TestResolve_MissingReference(t *testing.T) {
	env := map[string]string{"URL": "http://${HOST}:9000"}
	results := resolver.Resolve(env)
	r, ok := resultByKey(results, "URL")
	if !ok {
		t.Fatal("URL result not found")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "HOST" {
		t.Errorf("expected missing [HOST], got %v", r.Missing)
	}
	if r.Resolved != "http://<MISSING:HOST>:9000" {
		t.Errorf("unexpected resolved value: %q", r.Resolved)
	}
}

func TestResolve_MultipleRefs(t *testing.T) {
	env := map[string]string{
		"SCHEME": "https",
		"HOST":   "example.com",
		"URL":    "${SCHEME}://${HOST}/path",
	}
	results := resolver.Resolve(env)
	r, ok := resultByKey(results, "URL")
	if !ok {
		t.Fatal("URL result not found")
	}
	if r.Resolved != "https://example.com/path" {
		t.Errorf("unexpected resolved value: %q", r.Resolved)
	}
	sort.Strings(r.Refs)
	if len(r.Refs) != 2 || r.Refs[0] != "HOST" || r.Refs[1] != "SCHEME" {
		t.Errorf("unexpected refs: %v", r.Refs)
	}
}

func TestResolve_EmptyEnv(t *testing.T) {
	results := resolver.Resolve(map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
