# Project Scanner CLI

A simple open-source file scanner for source code projects.

Input a local folder path or a Git URL, and get a clear report of:

- Language distribution
- Per-language lines: code, comment, blank
- Largest files
- Project-level summary

## Why This Project

The tool is designed for easy onboarding:

- One command to analyze
- Text output for terminal usage
- JSON output for CI/automation
- HTML report for sharing and review

## Features (MVP)

- Scan source from local folder (`--path`).
- Scan source from public Git repo URL (`--url`).
- Analyze language usage and line counts.
- Analyze code/comment/blank breakdown.
- Show largest files (top N).
- Support output format `text`.
- Support output format `json`.
- Support output format `html`.
- Ignore built-in directories (`.git`, `node_modules`, `vendor`, `dist`, `build`, `coverage`).
- Read `.gitignore` patterns (basic support).
- Support custom ignore via `--ignore`.

## Quick Start

### 1. Clone and run

```bash
git clone <your-repo-url>
cd GITHUB-PST
go run ./cmd/scanner analyze --path .
```

### 2. Analyze a GitHub repository URL

```bash
go run ./cmd/scanner analyze --url https://github.com/golang/go
```

### 3. Export JSON

```bash
go run ./cmd/scanner analyze --path . --format json --out report.json --show-files
```

### 4. Export HTML report

```bash
go run ./cmd/scanner report --path . --out report.html
```

## CLI Usage

```bash
scanner analyze [--path . | --url <git-url>] [--format text|json|html] [--out report.html]
scanner report  [--path . | --url <git-url>] [--format html|json] [--out report.html]
scanner version
```

### Common options

- `--top 10` top N largest files
- `--workers 8` process files with 8 concurrent workers
- `--max-file-mb 5` skip files bigger than N MB
- `--ignore "*.min.js,docs/,tmp/"` extra ignore patterns
- `--show-files` include all file details in JSON output

## Suggested Distribution Strategy

For open-source usability, use this order:

1. GitHub Releases (main path): upload prebuilt binaries for Windows/macOS/Linux.
2. Package manager integration (next): Winget/Scoop/Homebrew.
3. GitHub Pages docs (optional): host docs and sample reports.
4. Source build path (developer path): `go build` and local testing.

## Build Binary

```bash
go build -o scanner ./cmd/scanner
```

## Project Structure

```text
cmd/scanner/main.go          # CLI entrypoint
internal/config              # Options and validation
internal/source              # Path/URL source preparation
internal/analyzer            # File walk, line counting, aggregation
internal/lang                # Language mapping and comment syntax
internal/output              # Text/JSON/HTML rendering
internal/model               # Report models
```

## Roadmap

- More accurate `.gitignore` pattern support
- Directory-level language hotspots
- Diff analysis between two scans
- GitHub Action summary comment for PR

## License

Apache-2.0
