package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github-pst/internal/config"
	"github-pst/internal/model"
)

func TestAnalyzeAggregatesByLanguage(t *testing.T) {
	root := t.TempDir()

	writeTestFile(t, root, "main.go", "package main\n\nfunc main() {\n\t// hello\n\tprintln(\"hi\")\n}\n")
	writeTestFile(t, root, "LICENSE", "line one\nline two\n")

	report, err := Analyze(model.SourceInfo{Kind: "path", Root: root}, config.Options{
		TopN:      10,
		Workers:   2,
		MaxFileMB: 5,
		ShowFiles: true,
	})
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if report.Summary.ScannedFiles != 2 {
		t.Fatalf("expected 2 scanned files, got %d", report.Summary.ScannedFiles)
	}
	if report.Summary.UnknownLanguage != 0 {
		t.Fatalf("expected 0 unknown language files, got %d", report.Summary.UnknownLanguage)
	}

	languages := map[string]model.LanguageStat{}
	for _, stat := range report.Languages {
		languages[stat.Language] = stat
	}

	if _, ok := languages["Go"]; !ok {
		t.Fatalf("expected Go language stat to exist")
	}
	if _, ok := languages["Text"]; !ok {
		t.Fatalf("expected Text language stat to exist")
	}
}

func TestAnalyzeHonorsIgnoreRulesAndShowFilesOption(t *testing.T) {
	root := t.TempDir()

	writeTestFile(t, root, "keep.go", "package main\nfunc main() {}\n")
	writeTestFile(t, root, "skip.tmp", "ignored\n")

	report, err := Analyze(model.SourceInfo{Kind: "path", Root: root}, config.Options{
		TopN:       10,
		Workers:    1,
		MaxFileMB:  5,
		ShowFiles:  false,
		IgnoreList: []string{"*.tmp"},
	})
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if report.Summary.ScannedFiles != 1 {
		t.Fatalf("expected 1 scanned file, got %d", report.Summary.ScannedFiles)
	}
	if report.Summary.SkippedFiles != 1 {
		t.Fatalf("expected 1 skipped file, got %d", report.Summary.SkippedFiles)
	}
	if report.Files != nil {
		t.Fatalf("expected Files to be nil when ShowFiles is false")
	}
}

func TestCountLinesGoSyntax(t *testing.T) {
	root := t.TempDir()
	path := writeTestFile(t, root, "sample.go", "package main\n\n// line comment\n/* block comment */\nfunc main() {}\n")

	code, comment, blank, total, err := countLines(path, "Go")
	if err != nil {
		t.Fatalf("countLines returned error: %v", err)
	}

	if code != 2 || comment != 2 || blank != 1 || total != 5 {
		t.Fatalf("unexpected counts: code=%d comment=%d blank=%d total=%d", code, comment, blank, total)
	}
}

func writeTestFile(t *testing.T, root, relPath, content string) string {
	t.Helper()
	path := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create dir for %s: %v", relPath, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", relPath, err)
	}
	return path
}
