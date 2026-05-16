// Package envnormaliser provides utilities for standardising the keys and
// values of parsed .env environments before comparison or reporting.
//
// It is intentionally non-destructive: all functions return new maps and never
// modify their inputs.
//
// Typical usage:
//
//	opts := envnormaliser.DefaultOptions()
//	normalised := envnormaliser.ApplyAll(loadedEnvs, opts)
//
// Options available:
//   - UppercaseKeys  – convert every key to UPPER_CASE (default: true)
//   - TrimValues     – strip leading/trailing whitespace from values (default: true)
//   - CollapseValues – collapse internal whitespace runs to a single space (default: false)
package envnormaliser
