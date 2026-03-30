package source

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github-pst/internal/model"
)

func Prepare(pathInput, urlInput string) (model.SourceInfo, func(), error) {
	if strings.TrimSpace(pathInput) != "" {
		abs, err := filepath.Abs(pathInput)
		if err != nil {
			return model.SourceInfo{}, nil, fmt.Errorf("resolve path: %w", err)
		}
		if _, err := os.Stat(abs); err != nil {
			return model.SourceInfo{}, nil, fmt.Errorf("invalid path: %w", err)
		}
		info := model.SourceInfo{Kind: "path", Input: pathInput, Root: abs}
		return info, func() {}, nil
	}

	if strings.TrimSpace(urlInput) != "" {
		tempDir, err := os.MkdirTemp("", "github-pst-*")
		if err != nil {
			return model.SourceInfo{}, nil, fmt.Errorf("create temp dir: %w", err)
		}

		cleanup := func() {
			_ = os.RemoveAll(tempDir)
		}

		cmd := exec.Command("git", "clone", "--depth", "1", urlInput, tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			cleanup()
			return model.SourceInfo{}, nil, fmt.Errorf("git clone failed: %w: %s", err, strings.TrimSpace(string(output)))
		}

		commit, _ := gitHeadCommit(tempDir)
		info := model.SourceInfo{Kind: "url", Input: urlInput, Root: tempDir, Commit: commit}
		return info, cleanup, nil
	}

	return model.SourceInfo{}, nil, errors.New("missing source input")
}

func gitHeadCommit(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
