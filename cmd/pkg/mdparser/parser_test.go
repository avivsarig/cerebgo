package mdparser

import (
	"path/filepath"
	"testing"
)

func TestParseMarkdownDoc(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{
			name:     "empty file",
			filepath: filepath.Join("testdata", "empty.md"),
			wantErr:  false,
		},
		{
			name:     "just frontmatter",
			filepath: filepath.Join("testdata", "frontmatter_only.md"),
			wantErr:  false,
		},
		{
			name:     "just content",
			filepath: filepath.Join("testdata", "content_only.md"),
			wantErr:  false,
		},
		{
			name:     "frontmatter and content",
			filepath: filepath.Join("testdata", "frontmatter_content.md"),
			wantErr:  false,
		},
		{
			name:     "content before frontmatter",
			filepath: filepath.Join("testdata", "invalid_content_before.md"),
			wantErr:  true,
		},
		{
			name:     "multiple frontmatter blocks",
			filepath: filepath.Join("testdata", "invalid_multiple_frontmatter.md"),
			wantErr:  true,
		},
		{
			name:     "incorrect dash count",
			filepath: filepath.Join("testdata", "invalid_dashes.md"),
			wantErr:  true,
		},
		{
			name:     "missing closing marker",
			filepath: filepath.Join("testdata", "invalid_unclosed.md"),
			wantErr:  true,
		},
		{
			name:     "whitespace before frontmatter",
			filepath: filepath.Join("testdata", "invalid_whitespace.md"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMarkdownDoc(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
