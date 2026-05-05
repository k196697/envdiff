package diff

import "sort"

// Status represents the comparison status of a key.
type Status string

const (
	StatusMatch      Status = "match"
	StatusMismatch   Status = "mismatch"
	StatusMissingInA Status = "missing_in_a"
	StatusMissingInB Status = "missing_in_b"
)

// Result holds the comparison result for a single key.
type Result struct {
	Key     string
	Status  Status
	ValueA  string
	ValueB  string
}

// Compare compares two env maps and returns a sorted slice of Results.
func Compare(a, b map[string]string) []Result {
	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	results := make([]Result, 0, len(keys))
	for k := range keys {
		va, inA := a[k]
		vb, inB := b[k]
		switch {
		case inA && inB && va == vb:
			results = append(results, Result{Key: k, Status: StatusMatch, ValueA: va, ValueB: vb})
		case inA && inB:
			results = append(results, Result{Key: k, Status: StatusMismatch, ValueA: va, ValueB: vb})
		case inA:
			results = append(results, Result{Key: k, Status: StatusMissingInB, ValueA: va})
		default:
			results = append(results, Result{Key: k, Status: StatusMissingInA, ValueB: vb})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// Summary returns counts of each status type.
func Summary(results []Result) (match, mismatch, missingA, missingB int) {
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			match++
		case StatusMismatch:
			mismatch++
		case StatusMissingInA:
			missingA++
		case StatusMissingInB:
			missingB++
		}
	}
	return
}
