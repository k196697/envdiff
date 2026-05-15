// Package resolver expands variable references within .env file values.
//
// Many .env files use interpolation to avoid repeating values:
//
//	HOST=localhost
//	DATABASE_URL=postgres://${HOST}/mydb
//
// Resolve walks every key in an environment map and substitutes
// ${VAR} or $VAR references with the corresponding value. Keys that
// cannot be resolved are flagged in the Result.Missing slice and
// replaced with a <MISSING:KEY> placeholder in the output so that
// downstream tooling can still process the string.
//
// WriteReport renders the results as plain text or JSON.
package resolver
