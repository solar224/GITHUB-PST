package analyzer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github-pst/internal/config"
	"github-pst/internal/lang"
	"github-pst/internal/model"
)

var defaultIgnoreContains = []string{
	".git",
	"node_modules",
	"vendor",
	"dist",
	"build",
	"coverage",
}

type ignoreMatcher struct {
	contains []string
	patterns []string
}

type scanTask struct {
	absPath string
	relPath string
	bytes   int64
}

type scanResult struct {
	stat model.FileStat
}

func Analyze(source model.SourceInfo, opts config.Options) (model.Report, error) {
	report := model.Report{
		GeneratedAt: time.Now().UTC(),
		Source:      source,
		Languages:   make([]model.LanguageStat, 0),
		Largest:     make([]model.FileStat, 0),
		Files:       make([]model.FileStat, 0),
	}

	matcher, err := newIgnoreMatcher(source.Root, opts.IgnoreList)
	if err != nil {
		return model.Report{}, err
	}

	langAgg := map[string]*model.LanguageStat{}
	maxFileBytes := opts.MaxFileMB * 1024 * 1024
	tasks := make([]scanTask, 0, 128)

	err = filepath.WalkDir(source.Root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		rel, err := filepath.Rel(source.Root, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			report.Summary.TotalDirs++
			return nil
		}

		if d.IsDir() {
			report.Summary.TotalDirs++
			if matcher.ShouldIgnore(rel, true) {
				return filepath.SkipDir
			}
			return nil
		}

		report.Summary.TotalFiles++
		if matcher.ShouldIgnore(rel, false) {
			report.Summary.SkippedFiles++
			report.Files = append(report.Files, model.FileStat{Path: rel, IsSkipped: true, SkipReason: "ignored by rule"})
			return nil
		}

		info, err := d.Info()
		if err != nil {
			report.Summary.SkippedFiles++
			report.Files = append(report.Files, model.FileStat{Path: rel, IsSkipped: true, SkipReason: "cannot read file info"})
			return nil
		}

		if info.Size() > maxFileBytes {
			report.Summary.SkippedFiles++
			report.Files = append(report.Files, model.FileStat{Path: rel, Bytes: info.Size(), IsSkipped: true, SkipReason: fmt.Sprintf("file too large (> %d MB)", opts.MaxFileMB)})
			return nil
		}

		tasks = append(tasks, scanTask{absPath: path, relPath: rel, bytes: info.Size()})

		return nil
	})
	if err != nil {
		return model.Report{}, fmt.Errorf("walk directory: %w", err)
	}

	results := processTasks(tasks, opts.Workers)
	for result := range results {
		report.Files = append(report.Files, result.stat)
		if result.stat.IsSkipped {
			report.Summary.SkippedFiles++
			continue
		}

		report.Summary.ScannedFiles++
		report.Summary.TotalCode += result.stat.Code
		report.Summary.TotalComment += result.stat.Comment
		report.Summary.TotalBlank += result.stat.Blank
		report.Summary.TotalLines += result.stat.Total
		report.Summary.TotalBytes += result.stat.Bytes
		if result.stat.Language == "Unknown" {
			report.Summary.UnknownLanguage++
		}

		if _, ok := langAgg[result.stat.Language]; !ok {
			langAgg[result.stat.Language] = &model.LanguageStat{Language: result.stat.Language}
		}
		agg := langAgg[result.stat.Language]
		agg.Files++
		agg.Code += result.stat.Code
		agg.Comment += result.stat.Comment
		agg.Blank += result.stat.Blank
		agg.Total += result.stat.Total
	}

	for _, stat := range langAgg {
		if report.Summary.TotalLines > 0 {
			stat.Percent = (float64(stat.Total) / float64(report.Summary.TotalLines)) * 100.0
		}
		report.Languages = append(report.Languages, *stat)
	}

	sort.Slice(report.Languages, func(i, j int) bool {
		return report.Languages[i].Total > report.Languages[j].Total
	})

	filesOnly := make([]model.FileStat, 0, len(report.Files))
	for _, f := range report.Files {
		if !f.IsSkipped {
			filesOnly = append(filesOnly, f)
		}
	}

	sort.Slice(filesOnly, func(i, j int) bool {
		return filesOnly[i].Bytes > filesOnly[j].Bytes
	})

	topN := opts.TopN
	if topN > len(filesOnly) {
		topN = len(filesOnly)
	}
	report.Largest = append(report.Largest, filesOnly[:topN]...)

	if !opts.ShowFiles {
		report.Files = nil
	}

	return report, nil
}

func processTasks(tasks []scanTask, workers int) <-chan scanResult {
	if workers <= 0 {
		workers = 1
	}
	if workers > len(tasks) && len(tasks) > 0 {
		workers = len(tasks)
	}

	taskCh := make(chan scanTask, workers*2+1)
	resultCh := make(chan scanResult, workers*2+1)

	if len(tasks) == 0 {
		close(resultCh)
		return resultCh
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				resultCh <- scanResult{stat: processTask(task)}
			}
		}()
	}

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
		wg.Wait()
		close(resultCh)
	}()

	return resultCh
}

func processTask(task scanTask) model.FileStat {
	if isBinaryFile(task.absPath) {
		return model.FileStat{Path: task.relPath, Bytes: task.bytes, IsSkipped: true, SkipReason: "binary file"}
	}

	language := lang.DetectLanguage(filepath.Base(task.absPath))
	code, comment, blank, total, err := countLines(task.absPath, language)
	if err != nil {
		return model.FileStat{Path: task.relPath, Bytes: task.bytes, Language: language, IsSkipped: true, SkipReason: "failed to parse lines"}
	}

	return model.FileStat{
		Path:     task.relPath,
		Language: language,
		Bytes:    task.bytes,
		Code:     code,
		Comment:  comment,
		Blank:    blank,
		Total:    total,
	}
}

func newIgnoreMatcher(root string, extra []string) (ignoreMatcher, error) {
	patterns := make([]string, 0)
	contains := append([]string{}, defaultIgnoreContains...)

	gitIgnore := filepath.Join(root, ".gitignore")
	if data, err := os.ReadFile(gitIgnore); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
				continue
			}
			patterns = append(patterns, strings.TrimPrefix(filepath.ToSlash(line), "./"))
		}
	}

	for _, e := range extra {
		e = strings.TrimSpace(filepath.ToSlash(e))
		if e == "" {
			continue
		}
		patterns = append(patterns, e)
	}

	return ignoreMatcher{contains: contains, patterns: patterns}, nil
}

func (m ignoreMatcher) ShouldIgnore(rel string, isDir bool) bool {
	normalized := strings.TrimPrefix(filepath.ToSlash(rel), "./")
	for _, item := range m.contains {
		if strings.Contains(normalized, "/"+item+"/") || strings.HasPrefix(normalized, item+"/") || normalized == item {
			return true
		}
	}

	for _, pattern := range m.patterns {
		p := strings.TrimSuffix(pattern, "/")
		if p == "" {
			continue
		}
		if ok, _ := filepath.Match(p, normalized); ok {
			return true
		}
		if strings.HasPrefix(normalized, p+"/") {
			return true
		}
		if isDir && normalized == p {
			return true
		}
	}

	return false
}

func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return true
	}
	buf = buf[:n]

	if bytes.IndexByte(buf, 0) != -1 {
		return true
	}
	return false
}

func countLines(path, language string) (int, int, int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer f.Close()

	syntax := lang.GetCommentSyntax(language)
	inBlock := false
	code := 0
	comment := 0
	blank := 0
	total := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		total++
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			blank++
			continue
		}

		if inBlock {
			comment++
			if syntax.BlockEnd != "" && strings.Contains(trimmed, syntax.BlockEnd) {
				inBlock = false
			}
			continue
		}

		isComment := false
		for _, prefix := range syntax.LinePrefixes {
			if strings.HasPrefix(trimmed, prefix) {
				comment++
				isComment = true
				break
			}
		}
		if isComment {
			continue
		}

		if syntax.SupportsBlock && syntax.BlockStart != "" && strings.HasPrefix(trimmed, syntax.BlockStart) {
			comment++
			if syntax.BlockEnd != "" && !strings.Contains(trimmed, syntax.BlockEnd) {
				inBlock = true
			}
			continue
		}

		code++
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, 0, err
	}

	return code, comment, blank, total, nil
}
