// Package envmatrix provides a cross-environment key/value matrix view.
//
// It aggregates multiple parsed .env files (each identified by a name such as
// "production", "staging", "development") into a single Matrix structure that
// makes it easy to see, at a glance, which keys are present or absent in each
// environment and what value each key holds.
//
// Basic usage:
//
//	envs := map[string]map[string]string{
//		"production":  {"DB_HOST": "prod-db", "API_KEY": "secret"},
//		"staging":     {"DB_HOST": "stage-db"},
//		"development": {"DB_HOST": "localhost", "DEBUG": "true"},
//	}
//
//	m := envmatrix.Build(envs)
//	envmatrix.Write(m, "text", os.Stdout)
//
// Keys and environment names are always sorted alphabetically for consistent,
// deterministic output. Missing keys are recorded in the Row.Absent slice and
// rendered as "<missing>" in text output.
package envmatrix
