package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Command string

const (
	CommandAnalyze Command = "analyze"
	CommandReport  Command = "report"
)

type Options struct {
	Command    Command
	Path       string
	URL        string
	Format     string
	Out        string
	TopN       int
	Workers    int
	ShowFiles  bool
	MaxFileMB  int64
	IgnoreList []string
}

func (o *Options) Validate() error {
	if strings.TrimSpace(o.Path) == "" && strings.TrimSpace(o.URL) == "" {
		return errors.New("you must provide either --path or --url")
	}
	if strings.TrimSpace(o.Path) != "" && strings.TrimSpace(o.URL) != "" {
		return errors.New("provide only one source: --path or --url")
	}
	o.Format = strings.ToLower(strings.TrimSpace(o.Format))
	if o.Format == "" {
		o.Format = "text"
	}
	if o.Format != "text" && o.Format != "json" && o.Format != "html" {
		return fmt.Errorf("invalid format %q, supported: text|json|html", o.Format)
	}
	if o.TopN <= 0 {
		o.TopN = 10
	}
	if o.Workers <= 0 {
		o.Workers = runtime.NumCPU()
		if o.Workers <= 0 {
			o.Workers = 1
		}
	}
	if o.MaxFileMB <= 0 {
		o.MaxFileMB = 5
	}
	if strings.TrimSpace(o.Path) != "" {
		o.Path = filepath.Clean(o.Path)
	}
	return nil
}
