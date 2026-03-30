package model

import "time"

type LanguageStat struct {
	Language string  `json:"language"`
	Files    int     `json:"files"`
	Code     int     `json:"code"`
	Comment  int     `json:"comment"`
	Blank    int     `json:"blank"`
	Total    int     `json:"total"`
	Percent  float64 `json:"percent"`
}

type FileStat struct {
	Path       string `json:"path"`
	Language   string `json:"language"`
	Bytes      int64  `json:"bytes"`
	Code       int    `json:"code"`
	Comment    int    `json:"comment"`
	Blank      int    `json:"blank"`
	Total      int    `json:"total"`
	IsSkipped  bool   `json:"isSkipped"`
	SkipReason string `json:"skipReason,omitempty"`
}

type SourceInfo struct {
	Kind   string `json:"kind"`
	Input  string `json:"input"`
	Root   string `json:"root"`
	Commit string `json:"commit,omitempty"`
}

type Summary struct {
	TotalFiles      int   `json:"totalFiles"`
	ScannedFiles    int   `json:"scannedFiles"`
	SkippedFiles    int   `json:"skippedFiles"`
	TotalDirs       int   `json:"totalDirs"`
	TotalCode       int   `json:"totalCode"`
	TotalComment    int   `json:"totalComment"`
	TotalBlank      int   `json:"totalBlank"`
	TotalLines      int   `json:"totalLines"`
	TotalBytes      int64 `json:"totalBytes"`
	UnknownLanguage int   `json:"unknownLanguage"`
}

type Report struct {
	GeneratedAt time.Time      `json:"generatedAt"`
	Source      SourceInfo     `json:"source"`
	Summary     Summary        `json:"summary"`
	Languages   []LanguageStat `json:"languages"`
	Largest     []FileStat     `json:"largestFiles"`
	Files       []FileStat     `json:"files"`
}
