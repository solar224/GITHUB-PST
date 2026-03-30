package lang

import "strings"

type CommentSyntax struct {
	LinePrefixes  []string
	BlockStart    string
	BlockEnd      string
	SupportsBlock bool
}

var extensionLanguage = map[string]string{
	".go":         "Go",
	".py":         "Python",
	".js":         "JavaScript",
	".mjs":        "JavaScript",
	".cjs":        "JavaScript",
	".ts":         "TypeScript",
	".tsx":        "TypeScript",
	".jsx":        "JavaScript",
	".java":       "Java",
	".rb":         "Ruby",
	".php":        "PHP",
	".rs":         "Rust",
	".c":          "C",
	".h":          "C",
	".cpp":        "C++",
	".cc":         "C++",
	".hpp":        "C++",
	".cs":         "C#",
	".swift":      "Swift",
	".kt":         "Kotlin",
	".kts":        "Kotlin",
	".scala":      "Scala",
	".sh":         "Shell",
	".bash":       "Shell",
	".zsh":        "Shell",
	".ps1":        "PowerShell",
	".sql":        "SQL",
	".html":       "HTML",
	".css":        "CSS",
	".scss":       "SCSS",
	".sass":       "SASS",
	".vue":        "Vue",
	".svelte":     "Svelte",
	".json":       "JSON",
	".xml":        "XML",
	".yaml":       "YAML",
	".yml":        "YAML",
	".toml":       "TOML",
	".md":         "Markdown",
	".dockerfile": "Dockerfile",
}

var languageCommentSyntax = map[string]CommentSyntax{
	"Go":         {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"JavaScript": {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"TypeScript": {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Java":       {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Rust":       {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"C":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"C++":        {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"C#":         {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Swift":      {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Kotlin":     {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Scala":      {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"PHP":        {LinePrefixes: []string{"//", "#"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Python":     {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Ruby":       {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Shell":      {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"PowerShell": {LinePrefixes: []string{"#"}, BlockStart: "<#", BlockEnd: "#>", SupportsBlock: true},
	"SQL":        {LinePrefixes: []string{"--"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"HTML":       {LinePrefixes: nil, BlockStart: "<!--", BlockEnd: "-->", SupportsBlock: true},
	"XML":        {LinePrefixes: nil, BlockStart: "<!--", BlockEnd: "-->", SupportsBlock: true},
	"CSS":        {LinePrefixes: nil, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"SCSS":       {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"SASS":       {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"YAML":       {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Dockerfile": {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"TOML":       {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Markdown":   {LinePrefixes: nil, SupportsBlock: false},
}

func DetectLanguage(fileName string) string {
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, "dockerfile") {
		return "Dockerfile"
	}
	for ext, language := range extensionLanguage {
		if strings.HasSuffix(lower, ext) {
			return language
		}
	}
	return "Unknown"
}

func GetCommentSyntax(language string) CommentSyntax {
	if syntax, ok := languageCommentSyntax[language]; ok {
		return syntax
	}
	return CommentSyntax{}
}
