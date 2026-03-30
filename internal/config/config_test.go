package config

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidateRequiresSource(t *testing.T) {
	opts := Options{}
	err := opts.Validate()
	if err == nil {
		t.Fatal("expected error when both path and url are missing")
	}
	if !strings.Contains(err.Error(), "either --path or --url") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsBothSources(t *testing.T) {
	opts := Options{Path: ".", URL: "https://github.com/example/repo"}
	err := opts.Validate()
	if err == nil {
		t.Fatal("expected error when both path and url are provided")
	}
	if !strings.Contains(err.Error(), "provide only one source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsInvalidFormat(t *testing.T) {
	opts := Options{Path: ".", Format: "xml"}
	err := opts.Validate()
	if err == nil {
		t.Fatal("expected invalid format error")
	}
	if !strings.Contains(err.Error(), "supported: text|json|html") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSetsDefaultsAndNormalizesPath(t *testing.T) {
	opts := Options{
		Path:      ".\\internal\\..\\internal\\config",
		Format:    "  ",
		TopN:      0,
		Workers:   0,
		MaxFileMB: 0,
	}

	if err := opts.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	if opts.Format != "text" {
		t.Fatalf("expected default format text, got %q", opts.Format)
	}
	if opts.TopN != 10 {
		t.Fatalf("expected default TopN=10, got %d", opts.TopN)
	}
	if opts.Workers <= 0 {
		t.Fatalf("expected Workers > 0, got %d", opts.Workers)
	}
	if opts.MaxFileMB != 5 {
		t.Fatalf("expected default MaxFileMB=5, got %d", opts.MaxFileMB)
	}
	if opts.Workers < 1 || opts.Workers > runtime.NumCPU() {
		t.Fatalf("expected Workers to be between 1 and NumCPU(%d), got %d", runtime.NumCPU(), opts.Workers)
	}

	expectedPath := filepath.Clean(".\\internal\\..\\internal\\config")
	if opts.Path != expectedPath {
		t.Fatalf("expected cleaned path %q, got %q", expectedPath, opts.Path)
	}
}

func TestValidateLowercasesFormat(t *testing.T) {
	opts := Options{Path: ".", Format: " JSON "}
	if err := opts.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if opts.Format != "json" {
		t.Fatalf("expected json format, got %q", opts.Format)
	}
}
