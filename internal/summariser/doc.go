// Package summariser aggregates per-file diff results produced by envdiff into
// a single top-level Summary value that can be rendered as plain text or JSON.
//
// Typical usage:
//
//	results := []summariser.FileResult{
//		{Name: "staging.env", Results: diffResults},
//	}
//	s := summariser.Build(results)
//	summariser.Write(os.Stdout, s, "text")
//
// The Summary.Healthy field is true only when every compared key matches
// across all file pairs — making it suitable as a quick pass/fail signal in
// CI pipelines.
package summariser
