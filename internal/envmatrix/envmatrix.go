// Package envmatrix builds a cross-environment key/value matrix,
// showing each key's value (or absence) across all loaded environments.
package envmatrix

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Row represents a single key across all environments.
type Row struct {
	Key    string            `json:"key"`
	Values map[string]string `json:"values"` // env name -> value (empty string means absent)
	Absent []string          `json:"absent,omitempty"`
}

// Matrix is the full cross-environment view.
type Matrix struct {
	Envs []string `json:"envs"`
	Rows []Row    `json:"rows"`
}

// Build constructs a Matrix from a map of env-name -> key/value pairs.
func Build(envs map[string]map[string]string) Matrix {
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]Row, 0, len(keys))
	for _, key := range keys {
		row := Row{
			Key:    key,
			Values: make(map[string]string, len(envNames)),
		}
		for _, env := range envNames {
			val, ok := envs[env][key]
			if ok {
				row.Values[env] = val
			} else {
				row.Absent = append(row.Absent, env)
			}
		}
		rows = append(rows, row)
	}

	return Matrix{Envs: envNames, Rows: rows}
}

// Write renders the matrix to w in the given format ("text" or "json").
func Write(m Matrix, format string, w io.Writer) error {
	switch strings.ToLower(format) {
	case "json":
		return writeJSON(m, w)
	default:
		return writeText(m, w)
	}
}

func writeText(m Matrix, w io.Writer) error {
	if len(m.Rows) == 0 {
		_, err := fmt.Fprintln(w, "no keys found")
		return err
	}
	// Header
	header := fmt.Sprintf("%-30s", "KEY")
	for _, env := range m.Envs {
		header += fmt.Sprintf("  %-20s", env)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, row := range m.Rows {
		line := fmt.Sprintf("%-30s", row.Key)
		for _, env := range m.Envs {
			val, ok := row.Values[env]
			if !ok {
				val = "<missing>"
			}
			if len(val) > 18 {
				val = val[:15] + "..."
			}
			line += fmt.Sprintf("  %-20s", val)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(m Matrix, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}
