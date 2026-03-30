package output

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github-pst/internal/model"
)

func RenderText(report model.Report) string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "Project Scanner Report\n")
	fmt.Fprintf(b, "Generated: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(b, "Source:    [%s] %s\n\n", report.Source.Kind, report.Source.Input)

	fmt.Fprintf(b, "Summary\n")
	fmt.Fprintf(b, "- Files:   %d scanned / %d skipped / %d total\n", report.Summary.ScannedFiles, report.Summary.SkippedFiles, report.Summary.TotalFiles)
	fmt.Fprintf(b, "- Dirs:    %d\n", report.Summary.TotalDirs)
	fmt.Fprintf(b, "- Lines:   %d code / %d comment / %d blank / %d total\n", report.Summary.TotalCode, report.Summary.TotalComment, report.Summary.TotalBlank, report.Summary.TotalLines)
	fmt.Fprintf(b, "- Bytes:   %d\n", report.Summary.TotalBytes)
	fmt.Fprintf(b, "- Unknown: %d files\n\n", report.Summary.UnknownLanguage)

	fmt.Fprintf(b, "Languages\n")
	for _, l := range report.Languages {
		fmt.Fprintf(b, "- %-14s %7d lines (%5.1f%%), files=%d\n", l.Language, l.Total, l.Percent, l.Files)
	}

	fmt.Fprintf(b, "\nLargest Files\n")
	for _, f := range report.Largest {
		fmt.Fprintf(b, "- %-40s %10d bytes (%s)\n", truncate(f.Path, 40), f.Bytes, f.Language)
	}

	return b.String()
}

func WriteJSON(report model.Report, out string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if out == "" {
		fmt.Println(string(data))
		return nil
	}
	if err := os.WriteFile(out, data, 0644); err != nil {
		return fmt.Errorf("write json output: %w", err)
	}
	return nil
}

func WriteHTML(report model.Report, out string) error {
	if out == "" {
		out = "report.html"
	}

	const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Project Scanner Report</title>
  <style>
    :root {
      --bg: #f6f4ef;
      --card: #fffdf7;
      --text: #1f1f1f;
      --accent: #0f766e;
      --muted: #6b7280;
      --line: #d8d5cc;
    }
    body { margin: 0; font-family: "Segoe UI", "Noto Sans", sans-serif; background: radial-gradient(circle at 5% 10%, #ece6d7, var(--bg)); color: var(--text); }
    .container { max-width: 1100px; margin: 0 auto; padding: 28px 18px 48px; }
    h1 { margin: 0 0 8px; letter-spacing: 0.5px; }
    .sub { color: var(--muted); margin-bottom: 20px; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 12px; margin-bottom: 16px; }
    .card { background: var(--card); border: 1px solid var(--line); border-radius: 12px; padding: 14px; box-shadow: 0 1px 0 #ece9df; }
    .value { font-size: 24px; font-weight: 700; color: var(--accent); }
    table { width: 100%; border-collapse: collapse; background: var(--card); border: 1px solid var(--line); border-radius: 12px; overflow: hidden; }
    th, td { padding: 10px; border-bottom: 1px solid var(--line); text-align: left; }
    th { background: #f2f0e8; }
    tr:last-child td { border-bottom: none; }
    .section { margin-top: 20px; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Project Scanner Report</h1>
    <div class="sub">Generated: {{ .GeneratedAt }} | Source: [{{ .Source.Kind }}] {{ .Source.Input }}</div>
    <div class="grid">
      <div class="card"><div>Total Files</div><div class="value">{{ .Summary.TotalFiles }}</div></div>
      <div class="card"><div>Scanned Files</div><div class="value">{{ .Summary.ScannedFiles }}</div></div>
      <div class="card"><div>Skipped Files</div><div class="value">{{ .Summary.SkippedFiles }}</div></div>
      <div class="card"><div>Total Lines</div><div class="value">{{ .Summary.TotalLines }}</div></div>
    </div>
    <div class="section">
      <h2>Language Distribution</h2>
      <table>
        <thead>
          <tr><th>Language</th><th>Files</th><th>Code</th><th>Comment</th><th>Blank</th><th>Total</th><th>Percent</th></tr>
        </thead>
        <tbody>
          {{ range .Languages }}
          <tr>
            <td>{{ .Language }}</td>
            <td>{{ .Files }}</td>
            <td>{{ .Code }}</td>
            <td>{{ .Comment }}</td>
            <td>{{ .Blank }}</td>
            <td>{{ .Total }}</td>
            <td>{{ printf "%.2f%%" .Percent }}</td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>

    <div class="section">
      <h2>Largest Files</h2>
      <table>
        <thead>
          <tr><th>Path</th><th>Language</th><th>Bytes</th><th>Total Lines</th></tr>
        </thead>
        <tbody>
          {{ range .Largest }}
          <tr>
            <td>{{ .Path }}</td>
            <td>{{ .Language }}</td>
            <td>{{ .Bytes }}</td>
            <td>{{ .Total }}</td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>
  </div>
</body>
</html>`

	t, err := template.New("report").Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse html template: %w", err)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create html output: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, report); err != nil {
		return fmt.Errorf("execute html template: %w", err)
	}
	return nil
}

func truncate(v string, max int) string {
	if len(v) <= max {
		return v
	}
	if max <= 3 {
		return v[:max]
	}
	return "..." + v[len(v)-max+3:]
}
