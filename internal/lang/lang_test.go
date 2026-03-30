package lang

import "testing"

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{name: "go file", fileName: "main.go", want: "Go"},
		{name: "typescript declaration", fileName: "index.d.ts", want: "TypeScript"},
		{name: "tsx uppercase", fileName: "App.TSX", want: "TypeScript"},
		{name: "dockerfile special", fileName: "Dockerfile", want: "Dockerfile"},
		{name: "go module special", fileName: "go.mod", want: "Go Module"},
		{name: "cmake special", fileName: "CMakeLists.txt", want: "CMake"},
		{name: "gitignore special", fileName: ".gitignore", want: "Shell"},
		{name: "markdown", fileName: "README.md", want: "Markdown"},
		{name: "txt file", fileName: "LICENSE.txt", want: "Text"},
		{name: "fallback text", fileName: "LICENSE", want: "Text"},
		{name: "empty fallback text", fileName: "", want: "Text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectLanguage(tt.fileName)
			if got != tt.want {
				t.Fatalf("DetectLanguage(%q) = %q, want %q", tt.fileName, got, tt.want)
			}
		})
	}
}

func TestGetCommentSyntax(t *testing.T) {
	goSyntax := GetCommentSyntax("Go")
	if len(goSyntax.LinePrefixes) == 0 || goSyntax.LinePrefixes[0] != "//" {
		t.Fatalf("unexpected Go line comment prefixes: %#v", goSyntax.LinePrefixes)
	}
	if !goSyntax.SupportsBlock || goSyntax.BlockStart != "/*" || goSyntax.BlockEnd != "*/" {
		t.Fatalf("unexpected Go block syntax: %#v", goSyntax)
	}

	unknownSyntax := GetCommentSyntax("NotALanguage")
	if len(unknownSyntax.LinePrefixes) != 0 || unknownSyntax.SupportsBlock {
		t.Fatalf("unknown language should have empty syntax, got %#v", unknownSyntax)
	}
}
