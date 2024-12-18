package mdparser

import (
	"path/filepath"
	"testing"
)

func TestParseMarkdownDoc(t *testing.T) {
	tests := []struct {
		name        string
		filepath    string
		wantErr     bool
		wantTitle   string
		wantContent string
	}{
		{
			name:      "empty file",
			filepath:  filepath.Join("testdata", "empty.md"),
			wantErr:   false,
			wantTitle: "empty",
		},
		{
			name:      "just frontmatter",
			filepath:  filepath.Join("testdata", "frontmatter_only.md"),
			wantErr:   false,
			wantTitle: "frontmatter_only",
		},
		{
			name:      "just content",
			filepath:  filepath.Join("testdata", "content_only.md"),
			wantErr:   false,
			wantTitle: "content_only",
		},
		{
			name:      "frontmatter and content",
			filepath:  filepath.Join("testdata", "frontmatter_content.md"),
			wantErr:   false,
			wantTitle: "frontmatter_content",
		},
		{
			name:      "content before frontmatter",
			filepath:  filepath.Join("testdata", "invalid_content_before.md"),
			wantErr:   true,
			wantTitle: "invalid_content_before", // Even with errors, title should be populated
		},
		{
			name:      "multiple frontmatter blocks",
			filepath:  filepath.Join("testdata", "invalid_multiple_frontmatter.md"),
			wantErr:   true,
			wantTitle: "invalid_multiple_frontmatter",
		},
		{
			name:      "incorrect dash count",
			filepath:  filepath.Join("testdata", "invalid_dashes.md"),
			wantErr:   true,
			wantTitle: "invalid_dashes",
		},
		{
			name:      "missing closing marker",
			filepath:  filepath.Join("testdata", "invalid_unclosed.md"),
			wantErr:   true,
			wantTitle: "invalid_unclosed",
		},
		{
			name:      "whitespace before frontmatter",
			filepath:  filepath.Join("testdata", "invalid_whitespace.md"),
			wantErr:   true,
			wantTitle: "invalid_whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := ParseMarkdownDoc(tt.filepath)

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, don't check other fields
			if tt.wantErr {
				return
			}

			// Verify the title is correctly populated from the filename
			if doc.Title != tt.wantTitle {
				t.Errorf("ParseMarkdownDoc() title = %v, want %v", doc.Title, tt.wantTitle)
			}
		})
	}
}
