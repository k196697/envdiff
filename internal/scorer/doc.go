// Package scorer provides a health-score calculation for envdiff results.
//
// It assigns a numeric score from 0 to 100 based on the proportion of
// matching, mismatched, and missing keys across compared .env files.
// Each status is weighted:
//
//	- Match    → full credit  (1.0)
//	- Mismatch → half credit  (0.5)
//	- Missing  → no credit    (0.0)
//
// The final score is rounded to two decimal places and mapped to a
// letter grade (A–F) for quick human-readable feedback.
package scorer
