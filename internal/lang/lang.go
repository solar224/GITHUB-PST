package lang

import (
	"sort"
	"strings"
)

type CommentSyntax struct {
	LinePrefixes  []string
	BlockStart    string
	BlockEnd      string
	SupportsBlock bool
}

var extensionLanguage = map[string]string{
	".adb":           "Ada",
	".ads":           "Ada",
	".asm":           "Assembly",
	".astro":         "Astro",
	".awk":           "Awk",
	".bash":          "Shell",
	".bat":           "Batch",
	".c":             "C",
	".cc":            "C++",
	".cjs":           "JavaScript",
	".clj":           "Clojure",
	".cljs":          "Clojure",
	".cls":           "TeX",
	".cmake":         "CMake",
	".coffee":        "CoffeeScript",
	".cpp":           "C++",
	".cs":            "C#",
	".css":           "CSS",
	".cxx":           "C++",
	".dart":          "Dart",
	".d.ts":          "TypeScript",
	".dockerfile":    "Dockerfile",
	".elm":           "Elm",
	".env":           "ENV",
	".erl":           "Erlang",
	".ex":            "Elixir",
	".exs":           "Elixir",
	".f90":           "Fortran",
	".f95":           "Fortran",
	".fs":            "F#",
	".fsi":           "F#",
	".fsx":           "F#",
	".gd":            "GDScript",
	".go":            "Go",
	".gomod":         "Go Module",
	".gosum":         "Go Module",
	".graphql":       "GraphQL",
	".groovy":        "Groovy",
	".h":             "C",
	".haml":          "Haml",
	".hbs":           "Handlebars",
	".hpp":           "C++",
	".hrl":           "Erlang",
	".hs":            "Haskell",
	".html":          "HTML",
	".hxx":           "C++",
	".ini":           "INI",
	".java":          "Java",
	".jl":            "Julia",
	".js":            "JavaScript",
	".json":          "JSON",
	".json5":         "JSON",
	".jsx":           "JavaScript",
	".kt":            "Kotlin",
	".kts":           "Kotlin",
	".less":          "LESS",
	".lhs":           "Haskell",
	".lisp":          "Lisp",
	".lua":           "Lua",
	".m":             "Objective-C",
	".make":          "Makefile",
	".md":            "Markdown",
	".mjml":          "MJML",
	".mjs":           "JavaScript",
	".mm":            "Objective-C++",
	".nim":           "Nim",
	".nix":           "Nix",
	".pas":           "Pascal",
	".php":           "PHP",
	".pl":            "Perl",
	".pm":            "Perl",
	".ps1":           "PowerShell",
	".psm1":          "PowerShell",
	".py":            "Python",
	".r":             "R",
	".rb":            "Ruby",
	".rs":            "Rust",
	".sass":          "SASS",
	".scala":         "Scala",
	".scm":           "Scheme",
	".scss":          "SCSS",
	".sh":            "Shell",
	".sol":           "Solidity",
	".sql":           "SQL",
	".svelte":        "Svelte",
	".swift":         "Swift",
	".tool-versions": "ENV",
	".tcl":           "Tcl",
	".tex":           "TeX",
	".tf":            "Terraform",
	".tfvars":        "Terraform",
	".txt":           "Text",
	".toml":          "TOML",
	".ts":            "TypeScript",
	".tsx":           "TypeScript",
	".twig":          "Twig",
	".vb":            "Visual Basic",
	".vue":           "Vue",
	".xml":           "XML",
	".yaml":          "YAML",
	".yml":           "YAML",
	".zig":           "Zig",
	".zsh":           "Shell",
}

var specialFileLanguage = map[string]string{
	"cmakelists.txt": "CMake",
	"containerfile":  "Dockerfile",
	"dockerfile":     "Dockerfile",
	"go.mod":         "Go Module",
	"go.sum":         "Go Module",
	"gemfile":        "Ruby",
	".gitattributes": "INI",
	".gitignore":     "Shell",
	".editorconfig":  "INI",
	"jenkinsfile":    "Groovy",
	"makefile":       "Makefile",
	"podfile":        "Ruby",
	"rakefile":       "Ruby",
	"vagrantfile":    "Ruby",
	"justfile":       "Makefile",
}

var sortedExtensions = sortedSuffixes(extensionLanguage)

var languageCommentSyntax = map[string]CommentSyntax{
	"Ada":           {LinePrefixes: []string{"--"}, SupportsBlock: false},
	"Assembly":      {LinePrefixes: []string{";", "#"}, SupportsBlock: false},
	"Astro":         {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Awk":           {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Batch":         {LinePrefixes: []string{"REM", "::"}, SupportsBlock: false},
	"C":             {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"C#":            {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"C++":           {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"CMake":         {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Clojure":       {LinePrefixes: []string{";"}, SupportsBlock: false},
	"CoffeeScript":  {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"CSS":           {LinePrefixes: nil, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Dart":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Dockerfile":    {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Elixir":        {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Elm":           {LinePrefixes: []string{"--"}, BlockStart: "{-", BlockEnd: "-}", SupportsBlock: true},
	"ENV":           {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Erlang":        {LinePrefixes: []string{"%"}, SupportsBlock: false},
	"F#":            {LinePrefixes: []string{"//"}, BlockStart: "(*", BlockEnd: "*)", SupportsBlock: true},
	"Fortran":       {LinePrefixes: []string{"!"}, SupportsBlock: false},
	"GDScript":      {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Go":            {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Go Module":     {LinePrefixes: []string{"//"}, SupportsBlock: false},
	"GraphQL":       {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Groovy":        {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Haml":          {LinePrefixes: []string{"-#", "/"}, SupportsBlock: false},
	"Handlebars":    {LinePrefixes: nil, BlockStart: "{{!--", BlockEnd: "--}}", SupportsBlock: true},
	"Haskell":       {LinePrefixes: []string{"--"}, BlockStart: "{-", BlockEnd: "-}", SupportsBlock: true},
	"HTML":          {LinePrefixes: nil, BlockStart: "<!--", BlockEnd: "-->", SupportsBlock: true},
	"INI":           {LinePrefixes: []string{";", "#"}, SupportsBlock: false},
	"Java":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"JavaScript":    {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Julia":         {LinePrefixes: []string{"#"}, BlockStart: "#=", BlockEnd: "=#", SupportsBlock: true},
	"Kotlin":        {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"LESS":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Lisp":          {LinePrefixes: []string{";"}, BlockStart: "#|", BlockEnd: "|#", SupportsBlock: true},
	"Lua":           {LinePrefixes: []string{"--"}, BlockStart: "--[[", BlockEnd: "]]", SupportsBlock: true},
	"Makefile":      {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Markdown":      {LinePrefixes: nil, SupportsBlock: false},
	"MJML":          {LinePrefixes: nil, BlockStart: "<!--", BlockEnd: "-->", SupportsBlock: true},
	"Nim":           {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Nix":           {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Objective-C":   {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Objective-C++": {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Pascal":        {LinePrefixes: []string{"//"}, BlockStart: "{", BlockEnd: "}", SupportsBlock: true},
	"Perl":          {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"PHP":           {LinePrefixes: []string{"//", "#"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"PowerShell":    {LinePrefixes: []string{"#"}, BlockStart: "<#", BlockEnd: "#>", SupportsBlock: true},
	"Python":        {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"R":             {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Ruby":          {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Rust":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"SASS":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Scala":         {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Scheme":        {LinePrefixes: []string{";"}, SupportsBlock: false},
	"SCSS":          {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Shell":         {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Solidity":      {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"SQL":           {LinePrefixes: []string{"--"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Svelte":        {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Swift":         {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Tcl":           {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"TeX":           {LinePrefixes: []string{"%"}, SupportsBlock: false},
	"Terraform":     {LinePrefixes: []string{"#", "//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"TOML":          {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Twig":          {LinePrefixes: nil, BlockStart: "{#", BlockEnd: "#}", SupportsBlock: true},
	"TypeScript":    {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"Visual Basic":  {LinePrefixes: []string{"'", "REM"}, SupportsBlock: false},
	"Vue":           {LinePrefixes: []string{"//"}, BlockStart: "/*", BlockEnd: "*/", SupportsBlock: true},
	"XML":           {LinePrefixes: nil, BlockStart: "<!--", BlockEnd: "-->", SupportsBlock: true},
	"YAML":          {LinePrefixes: []string{"#"}, SupportsBlock: false},
	"Zig":           {LinePrefixes: []string{"//"}, SupportsBlock: false},
}

func DetectLanguage(fileName string) string {
	lower := strings.ToLower(strings.TrimSpace(fileName))
	if lower == "" {
		return "Text"
	}

	if language, ok := specialFileLanguage[lower]; ok {
		return language
	}

	for _, ext := range sortedExtensions {
		language := extensionLanguage[ext]
		if strings.HasSuffix(lower, ext) {
			return language
		}
	}

	return "Text"
}

func sortedSuffixes(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) == len(keys[j]) {
			return keys[i] < keys[j]
		}
		return len(keys[i]) > len(keys[j])
	})

	return keys
}

func GetCommentSyntax(language string) CommentSyntax {
	if syntax, ok := languageCommentSyntax[language]; ok {
		return syntax
	}
	return CommentSyntax{}
}
