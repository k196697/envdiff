package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its name and parsed key-value pairs.
type EnvFile struct {
	Name string
	Path string
	Env  map[string]string
}

// LoadAll loads multiple .env files by path, returning a slice of EnvFile.
// It returns an error if any file cannot be parsed.
func LoadAll(paths []string) ([]EnvFile, error) {
	files := make([]EnvFile, 0, len(paths))
	for _, p := range paths {
		env, err := parser.ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("loader: failed to parse %q: %w", p, err)
		}
		files = append(files, EnvFile{
			Name: filepath.Base(p),
			Path: p,
			Env:  env,
		})
	}
	return files, nil
}

// LoadDir scans a directory for files matching the given glob pattern (e.g. "*.env", ".env*")
// and loads all matching files.
func LoadDir(dir string, pattern string) ([]EnvFile, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil, fmt.Errorf("loader: glob error: %w", err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("loader: no files matched pattern %q in %q", pattern, dir)
	}
	var paths []string
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil || info.IsDir() {
			continue
		}
		paths = append(paths, m)
	}
	return LoadAll(paths)
}
