// Package classifier categorises diff results by severity level.
// Each result is assigned a severity (info, warning, critical) based
// on its status and optional user-supplied rules.
package classifier

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Severity represents how serious a diff result is.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Result pairs a diff result with its assigned severity.
type Result struct {
	Diff     diff.Result
	Severity Severity
}

// Options controls classification behaviour.
type Options struct {
	// CriticalPrefixes marks any key with these prefixes as critical.
	CriticalPrefixes []string
	// WarningPrefixes marks any key with these prefixes as warning (unless
	// already critical).
	WarningPrefixes []string
}

// Classify assigns a Severity to each diff.Result.
func Classify(results []diff.Result, opts Options) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Diff:     r,
			Severity: assign(r, opts),
		})
	}
	return out
}

func assign(r diff.Result, opts Options) Severity {
	key := strings.ToUpper(r.Key)

	for _, p := range opts.CriticalPrefixes {
		if strings.HasPrefix(key, strings.ToUpper(p)) {
			return SeverityCritical
		}
	}

	for _, p := range opts.WarningPrefixes {
		if strings.HasPrefix(key, strings.ToUpper(p)) {
			return SeverityWarning
		}
	}

	switch r.Status {
	case diff.StatusMissing:
		return SeverityWarning
	case diff.StatusMismatch:
		return SeverityWarning
	default:
		return SeverityInfo
	}
}
