package loader_test

import (
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func TestValidate_AllValid(t *testing.T) {
	files := []loader.EnvFile{
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{"KEY": "val"}},
		{Name: ".env.dev", Path: "/a/.env.dev", Env: map[string]string{"KEY": "other"}},
	}
	if err := loader.Validate(files); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_DuplicatePath(t *testing.T) {
	files := []loader.EnvFile{
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{"KEY": "val"}},
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{"KEY": "val"}},
	}
	err := loader.Validate(files)
	if err == nil {
		t.Fatal("expected error for duplicate path, got nil")
	}
	ve, ok := err.(*loader.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d: %v", len(ve.Issues), ve.Issues)
	}
}

func TestValidate_EmptyEnv(t *testing.T) {
	files := []loader.EnvFile{
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{}},
	}
	err := loader.Validate(files)
	if err == nil {
		t.Fatal("expected error for empty env, got nil")
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	files := []loader.EnvFile{
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{}},
		{Name: ".env.prod", Path: "/a/.env.prod", Env: map[string]string{}},
	}
	err := loader.Validate(files)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	ve := err.(*loader.ValidationError)
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}
