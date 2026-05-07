// Package patcher provides functionality for patching .env files by appending
// keys that are present in a reference environment but missing from the target.
//
// It integrates with the diff package to consume comparison results and writes
// only the keys whose status is diff.Missing. Existing keys are never modified.
//
// Usage:
//
//	results := diff.Compare(reference, target)
//	res, err := patcher.Apply(".env.local", results, patcher.Options{
//		DryRun:      false,
//		Placeholder: "CHANGEME",
//	})
//
// When DryRun is true the file is not modified; the returned Result still
// reports which keys would have been added, making it safe to use in
// preview or CI reporting workflows.
package patcher
