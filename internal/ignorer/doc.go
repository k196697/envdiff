// Package ignorer loads an optional .envdiffignore file and provides
// key-matching and filtering helpers so that certain keys can be
// universally excluded from all diff comparisons.
//
// File format:
//
//	# Lines beginning with '#' are comments and are ignored.
//	# Blank lines are also ignored.
//	# Each non-blank, non-comment line is treated as a key name.
//	# Key matching is case-insensitive.
//
// Example .envdiffignore:
//
//	# Internal secrets — never compare these
//	AWS_SECRET_ACCESS_KEY
//	DB_PASSWORD
//	JWT_SECRET
package ignorer
