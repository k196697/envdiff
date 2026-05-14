package renamer

// Similarity exposes the internal similarity function for testing and external
// use without needing to construct full diff.Result slices.
func Similarity(a, b string) float64 {
	return similarity(a, b)
}
