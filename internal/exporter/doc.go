// Package exporter provides functionality for persisting envdiff results
// to disk in a variety of formats.
//
// Supported formats:
//   - text     — plain key/status lines, human-readable
//   - json     — structured JSON suitable for CI pipelines
//   - markdown — GitHub-flavoured markdown table
//
// Usage:
//
//	err := exporter.Export(results, st, exporter.Options{
//		OutputPath: "report.md",
//		Format:     exporter.FormatMarkdown,
//		Stats:      true,
//	})
package exporter
