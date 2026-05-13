package grouper_test

import (
	"fmt"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func BenchmarkByPrefix_Large(b *testing.B) {
	const n = 2000
	results := make([]diff.Result, n)
	prefixes := []string{"DB", "APP", "AWS", "GCP", "REDIS", "KAFKA", "AUTH", "SMTP"}
	for i := 0; i < n; i++ {
		prefix := prefixes[i%len(prefixes)]
		results[i] = diff.Result{
			Key:    fmt.Sprintf("%s_KEY_%d", prefix, i),
			Status: diff.Match,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.ByPrefix(results, grouper.Options{})
	}
}

func BenchmarkByPrefix_Depth2(b *testing.B) {
	const n = 1000
	results := make([]diff.Result, n)
	for i := 0; i < n; i++ {
		results[i] = diff.Result{
			Key:    fmt.Sprintf("AWS_S3_KEY_%d", i),
			Status: diff.Mismatch,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grouper.ByPrefix(results, grouper.Options{Depth: 2})
	}
}
