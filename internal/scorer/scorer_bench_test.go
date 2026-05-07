package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

// BenchmarkCompute_Large measures scorer throughput on a large result set.
func BenchmarkCompute_Large(b *testing.B) {
	statuses := []string{
		diff.StatusMatch,
		diff.StatusMismatch,
		diff.StatusMissingInA,
		diff.StatusMissingInB,
	}

	results := make([]diff.Result, 1000)
	for i := range results {
		results[i] = diff.Result{
			Key:    "KEY",
			Status: statuses[i%len(statuses)],
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scorer.Compute(results)
	}
}

// BenchmarkCompute_AllMatch is a best-case scenario benchmark.
func BenchmarkCompute_AllMatch(b *testing.B) {
	results := make([]diff.Result, 500)
	for i := range results {
		results[i] = diff.Result{Key: "KEY", Status: diff.StatusMatch}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scorer.Compute(results)
	}
}
