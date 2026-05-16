// Package trimmer detects values in diff results that contain leading or
// trailing whitespace characters (spaces and tabs).
//
// Such values are syntactically valid in .env files but frequently cause
// subtle bugs when the surrounding application trims or does not trim input
// before comparison.
//
// Usage:
//
//	issues := trimmer.Detect(results)
//	_ = trimmer.WriteReport(os.Stdout, issues, "text")
//
// Detect returns an []Issue slice sorted by file then key. WriteReport
// supports "text" and "json" output formats; any unknown format falls back
// to plain text.
package trimmer
