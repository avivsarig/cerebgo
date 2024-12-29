package mdparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestWriteMarkdownDoc(t *testing.T) {
	// Define test cases using table-driven approach
	tests := []struct {
		name    string
		fm      mdparser.Frontmatter
		content string
		setup   func(t *testing.T) string // returns file path
		wantErr bool
		want    string // expected file content
	}{
		{
			name: "simple document",
			fm: mdparser.Frontmatter{
				"title":  "Test Doc",
				"author": "Tester",
			},
			content: "Hello, World!",
			setup: func(t *testing.T) string {
				return filepath.Join(testutil.CreateTestDirectory(t), "test.md")
			},
			wantErr: false,
			want: `---
author: Tester
title: Test Doc
---

Hello, World!`,
		},
		{
			name: "complex frontmatter",
			fm: mdparser.Frontmatter{
				"nested": map[string]interface{}{
					"key": "value",
				},
				"list": []string{"item1", "item2"},
			},
			content: "Complex content",
			setup: func(t *testing.T) string {
				return filepath.Join(testutil.CreateTestDirectory(t), "complex.md")
			},
			wantErr: false,
			want: `---
list:
    - item1
    - item2
nested:
    key: value
---

Complex content`,
		},
		{
			name: "write error - invalid path",
			fm: mdparser.Frontmatter{
				"title": "Test",
			},
			content: "Content",
			setup: func(t *testing.T) string {
				return "/nonexistent/directory/file.md"
			},
			wantErr: true,
		},
		{
			name: "write error - invalid permissions",
			fm: mdparser.Frontmatter{
				"title": "Test",
			},
			content: "Content",
			setup: func(t *testing.T) string {
				dir := testutil.CreateTestDirectory(t)
				path := filepath.Join(dir, "readonly.md")
				// Create read-only directory
				if err := os.Chmod(dir, 0555); err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr: true,
		},
		{
			name:    "empty frontmatter",
			fm:      mdparser.Frontmatter{},
			content: "Just content",
			setup: func(t *testing.T) string {
				return filepath.Join(testutil.CreateTestDirectory(t), "empty.md")
			},
			want: `---
---

Just content`,
		},
		{
			name: "unicode content",
			fm: mdparser.Frontmatter{
				"title":  "שלום",
				"author": "测试",
			},
			content: "Hello مرحبا こんにちは",
			setup: func(t *testing.T) string {
				return filepath.Join(testutil.CreateTestDirectory(t), "unicode.md")
			},
			want: `---
author: 测试
title: שלום
---

Hello مرحبا こんにちは`,
		},
		{
			name: "invalid frontmatter - func value",
			fm: mdparser.Frontmatter{
				"func": (func())(nil), // Using nil func for deterministic error
			},
			content: "Content",
			setup: func(t *testing.T) string {
				return filepath.Join(testutil.CreateTestDirectory(t), "invalid.md")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)

			err := mdparser.WriteMarkdownDoc(tt.fm, tt.content, path)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteMarkdownDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				testutil.AssertFileContent(t, path, tt.want)
			}
		})
	}
}
