package config_test

import (
	"flag"
	"testing"

	"github.com/user/envdiff/internal/config"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestParse_BasicFiles(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{".env", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestParse_DirFlag(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-dir", "/envs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Dir != "/envs" {
		t.Errorf("expected dir /envs, got %s", cfg.Dir)
	}
}

func TestParse_FormatJSON(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected json format, got %s", cfg.Format)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-format", "yaml"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParse_ExcludeKeys(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-exclude", "SECRET,TOKEN"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Exclude) != 2 {
		t.Errorf("expected 2 excluded keys, got %d", len(cfg.Exclude))
	}
}

func TestParse_ExportFlags(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{
		"-output", "report.md",
		"-export-format", "markdown",
		"-export-stats",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputPath != "report.md" {
		t.Errorf("expected output report.md, got %s", cfg.OutputPath)
	}
	if cfg.ExportFmt != "markdown" {
		t.Errorf("expected markdown export format, got %s", cfg.ExportFmt)
	}
	if !cfg.ExportStats {
		t.Error("expected ExportStats to be true")
	}
}

func TestParse_InvalidExportFormat(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-export-format", "csv"})
	if err == nil {
		t.Fatal("expected error for invalid export-format")
	}
}

func TestParse_InvalidSortBy(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-sort", "value"})
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}
}
