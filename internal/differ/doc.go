// Package differ orchestrates multi-file environment comparisons.
//
// It compares every loaded environment file against a single designated
// baseline, returning a slice of FileResult values — one per non-baseline
// file — that can be fed into the reporter, formatter, or exporter packages.
//
// # Basic usage
//
//	results, err := differ.RunFiles([]string{".env", ".env.staging", ".env.prod"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	summaries := differ.Summarise(results)
//	for _, s := range summaries {
//		fmt.Printf("%s vs %s: %d mismatches\n", s.Baseline, s.Target, s.Mismatch)
//	}
//
// The first path passed to RunFiles is always treated as the baseline.
// Use RunAll when you have already parsed the env maps yourself.
package differ
