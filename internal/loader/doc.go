// Package loader provides utilities for discovering and loading multiple .env
// files from explicit paths or by scanning a directory with a glob pattern.
//
// Each loaded file is represented as an EnvFile containing its base name, full
// path, and the parsed key-value map produced by the parser package.
//
// Typical usage:
//
//	files, err := loader.LoadAll([]string{".env.prod", ".env.staging"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Or using directory scanning:
//
//	files, err := loader.LoadDir(".", ".env.*")
package loader
