package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Match},
		{Key: "DB_PORT", Status: diff.Mismatch},
		{Key: "APP_ENV", Status: diff.Match},
		{Key: "APP_DEBUG", Status: diff.MissingInB},
		{Key: "STANDALONE", Status: diff.Match},
	}
}

func TestByPrefix_DefaultDepth(t *testing.T) {
	groups := grouper.ByPrefix(sampleResults(), grouper.Options{})

	prefixMap := make(map[string]int)
	for _, g := range groups {
		prefixMap[g.Prefix] = len(g.Results)
	}

	if prefixMap["DB"] != 2 {
		t.Errorf("expected 2 DB results, got %d", prefixMap["DB"])
	}
	if prefixMap["APP"] != 2 {
		t.Errorf("expected 2 APP results, got %d", prefixMap["APP"])
	}
	if prefixMap["(other)"] != 1 {
		t.Errorf("expected 1 (other) result, got %d", prefixMap["(other)"])
	}
}

func TestByPrefix_Sorted(t *testing.T) {
	groups := grouper.ByPrefix(sampleResults(), grouper.Options{})
	for i := 1; i < len(groups); i++ {
		if groups[i].Prefix < groups[i-1].Prefix {
			t.Errorf("groups not sorted: %q before %q", groups[i-1].Prefix, groups[i].Prefix)
		}
	}
}

func TestByPrefix_CustomSeparator(t *testing.T) {
	results := []diff.Result{
		{Key: "aws.region", Status: diff.Match},
		{Key: "aws.secret", Status: diff.Mismatch},
		{Key: "gcp.project", Status: diff.Match},
	}
	groups := grouper.ByPrefix(results, grouper.Options{Separator: "."})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestByPrefix_Empty(t *testing.T) {
	groups := grouper.ByPrefix(nil, grouper.Options{})
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestByPrefix_Depth2(t *testing.T) {
	results := []diff.Result{
		{Key: "AWS_S3_BUCKET", Status: diff.Match},
		{Key: "AWS_S3_REGION", Status: diff.Match},
		{Key: "AWS_EC2_TYPE", Status: diff.Mismatch},
	}
	groups := grouper.ByPrefix(results, grouper.Options{Depth: 2})
	prefixMap := make(map[string]int)
	for _, g := range groups {
		prefixMap[g.Prefix] = len(g.Results)
	}
	if prefixMap["AWS_S3"] != 2 {
		t.Errorf("expected 2 AWS_S3 results, got %d", prefixMap["AWS_S3"])
	}
	if prefixMap["AWS_EC2"] != 1 {
		t.Errorf("expected 1 AWS_EC2 result, got %d", prefixMap["AWS_EC2"])
	}
}
