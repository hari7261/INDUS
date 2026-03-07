package main

import (
	"path/filepath"
	"testing"
)

func TestParseBootstrapArgsWithoutBootstrapFlags(t *testing.T) {
	cfg, err := parseBootstrapArgs([]string{"ind", "sys", "stats"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CWD != "" {
		t.Fatalf("expected empty cwd, got %q", cfg.CWD)
	}
	if len(cfg.Tokens) != 3 {
		t.Fatalf("unexpected token count: %d", len(cfg.Tokens))
	}
}

func TestParseBootstrapArgsWithCWD(t *testing.T) {
	cfg, err := parseBootstrapArgs([]string{"--cwd", ".", "ind", "version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CWD != "." {
		t.Fatalf("expected cwd override '.', got %q", cfg.CWD)
	}
	if len(cfg.Tokens) != 2 {
		t.Fatalf("unexpected token count: %d", len(cfg.Tokens))
	}
}

func TestParseBootstrapArgsRequiresCWDValue(t *testing.T) {
	_, err := parseBootstrapArgs([]string{"--cwd"})
	if err == nil {
		t.Fatal("expected parse error for missing --cwd value")
	}
}

func TestResolveWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	resolved, err := resolveWorkingDirectory(root)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	abs, _ := filepath.Abs(root)
	if resolved != abs {
		t.Fatalf("unexpected resolved path: %q (expected %q)", resolved, abs)
	}
}

func TestResolveWorkingDirectoryRejectsMissingDir(t *testing.T) {
	_, err := resolveWorkingDirectory(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected resolve error for missing directory")
	}
}
