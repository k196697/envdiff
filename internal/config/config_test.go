package config

import (
	"testing"
)

func TestParse_BasicFiles(t *testing.T) {
	opts, err := Parse([]string{".env.dev", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(opts.Files))
	}
	if opts.Format != FormatText {
		t.Errorf("expected default format text, got %q", opts.Format)
	}
	if opts.Quiet {
		t.Error("expected quiet=false by default")
	}
}

func TestParse_DirFlag(t *testing.T) {
	opts, err := Parse([]string{"--dir", "./envs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Dir != "./envs" {
		t.Errorf("expected dir ./envs, got %q", opts.Dir)
	}
	if len(opts.Files) != 0 {
		t.Errorf("expected no positional files, got %d", len(opts.Files))
	}
}

func TestParse_FormatJSON(t *testing.T) {
	opts, err := Parse([]string{"--format", "json", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != FormatJSON {
		t.Errorf("expected json format, got %q", opts.Format)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]string{"--format", "yaml", ".env.a", ".env.b"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParse_ExcludeKeys(t *testing.T) {
	opts, err := Parse([]string{"--exclude", "SECRET, TOKEN , DEBUG", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.ExcludeKeys) != 3 {
		t.Errorf("expected 3 exclude keys, got %d", len(opts.ExcludeKeys))
	}
	if opts.ExcludeKeys[1] != "TOKEN" {
		t.Errorf("expected trimmed key TOKEN, got %q", opts.ExcludeKeys[1])
	}
}

func TestParse_TooFewFiles(t *testing.T) {
	_, err := Parse([]string{".env.only"})
	if err == nil {
		t.Fatal("expected error when fewer than 2 files and no --dir")
	}
}

func TestParse_DirAndFilesConflict(t *testing.T) {
	_, err := Parse([]string{"--dir", "./envs", ".env.extra"})
	if err == nil {
		t.Fatal("expected error when both --dir and positional files are provided")
	}
}

func TestParse_QuietAndPrefix(t *testing.T) {
	opts, err := Parse([]string{"--quiet", "--prefix", "APP_", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Quiet {
		t.Error("expected quiet=true")
	}
	if opts.Prefix != "APP_" {
		t.Errorf("expected prefix APP_, got %q", opts.Prefix)
	}
}
