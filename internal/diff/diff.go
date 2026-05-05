package diff

// KeyStatus represents the comparison result for a single key.
type KeyStatus int

const (
	StatusMatch    KeyStatus = iota // key exists in both with same value
	StatusMismatch                  // key exists in both but values differ
	StatusMissingB                  // key exists in A but not in B
	StatusMissingA                  // key exists in B but not in A
)

// Result holds the comparison result for a single key across two env files.
type Result struct {
	Key      string
	ValueA   string
	ValueB   string
	Status   KeyStatus
}

// Compare takes two parsed env maps (fileA, fileB) and returns a slice of
// Result entries describing every key found in either file.
func Compare(fileA, fileB map[string]string) []Result {
	seen := make(map[string]bool)
	var results []Result

	for k, vA := range fileA {
		seen[k] = true
		vB, ok := fileB[k]
		switch {
		case !ok:
			results = append(results, Result{Key: k, ValueA: vA, Status: StatusMissingB})
		case vA == vB:
			results = append(results, Result{Key: k, ValueA: vA, ValueB: vB, Status: StatusMatch})
		default:
			results = append(results, Result{Key: k, ValueA: vA, ValueB: vB, Status: StatusMismatch})
		}
	}

	for k, vB := range fileB {
		if !seen[k] {
			results = append(results, Result{Key: k, ValueB: vB, Status: StatusMissingA})
		}
	}

	return results
}

// Summary returns counts of each status across the results.
func Summary(results []Result) (match, mismatch, missingB, missingA int) {
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			match++
		case StatusMismatch:
			mismatch++
		case StatusMissingB:
			missingB++
		case StatusMissingA:
			missingA++
		}
	}
	return
}
