package mdparser_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestParseMarkdownDocFileHandling(t *testing.T) {
	// Create temporary test directory
	testDir := testutil.CreateTestDirectory(t)

	tests := []struct {
		name        string
		content     string
		want        mdparser.MarkdownDocument
		wantErr     bool
		errContains string
	}{
		{
			name:    "empty file",
			content: "",
			want: mdparser.MarkdownDocument{
				Title: "test",
			},
		},
		{
			name:    "content only, no frontmatter",
			content: "# Hello World",
			want: mdparser.MarkdownDocument{
				Title:   "test",
				Content: "# Hello World",
			},
		},
		{
			name: "valid frontmatter and content",
			content: `---
title: Test
date: 2024-01-01
---
# Content here`,
			want: mdparser.MarkdownDocument{
				Title: "test",
				Frontmatter: map[string]interface{}{
					"title": "Test",
					"date":  "2024-01-01",
				},
				Content: "# Content here",
			},
		},
		{
			name: "frontmatter only",
			content: `---
key: value
---`,
			want: mdparser.MarkdownDocument{
				Title: "test",
				Frontmatter: map[string]interface{}{
					"key": "value",
				},
				Content: "",
			},
		},
		{
			name: "whitespace before frontmatter",
			content: ` ---
key: value
---`,
			wantErr:     true,
			errContains: "whitespace before frontmatter",
		},
		{
			name: "invalid frontmatter marker",
			content: `----
key: value
---`,
			wantErr:     true,
			errContains: "incorrect frontmatter markers",
		},
		{
			name: "unclosed frontmatter",
			content: `---
key: value`,
			wantErr:     true,
			errContains: "unclosed frontmatter",
		},
		{
			name: "multiple frontmatter blocks",
			content: `---
key1: value1
---
content
---
key2: value2
---`,
			wantErr:     true,
			errContains: "multiple frontmatter blocks",
		},
		{
			name: "invalid YAML in frontmatter",
			content: `---
key: : invalid : yaml :
---`,
			wantErr:     true,
			errContains: "invalid frontmatter YAML",
		},
		{
			name:        "dashes in content without frontmatter",
			content:     "Some content\n---\nMore content",
			wantErr:     true,
			errContains: "content before frontmatter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filename := "test.md"
			filepath := filepath.Join(testDir, filename)
			err := testutil.CreateTestFile(t, testDir, filename, tt.content)
			if err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			// Run test
			got, err := mdparser.ParseMarkdownDoc(filepath)

			// Verify error cases
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseMarkdownDoc() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want it to contain %v", err, tt.errContains)
				}
				return
			}

			// Verify success cases
			if err != nil {
				t.Errorf("ParseMarkdownDoc() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMarkdownDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMarkdownDoc_FileErrors(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		wantErr     bool
		errContains string
	}{
		{
			name: "non-existent file",
			setup: func(t *testing.T) string {
				return "/nonexistent/path/file.md"
			},
			wantErr:     true,
			errContains: "failed to read file",
		},
		{
			name: "permission denied",
			setup: func(t *testing.T) string {
				dir := testutil.CreateTestDirectory(t)
				path := filepath.Join(dir, "noperm.md")
				err := testutil.CreateTestFile(t, dir, "noperm.md", "content")
				if err != nil {
					t.Fatal(err)
				}
				// Remove read permissions
				if err := os.Chmod(path, 0000); err != nil {
					t.Fatal(err)
				}
				return path
			},
			wantErr:     true,
			errContains: "failed to read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filepath := tt.setup(t)
			_, err := mdparser.ParseMarkdownDoc(filepath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseMarkdownDoc() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want it to contain %v", err, tt.errContains)
				}
				return
			}
		})
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(substr) > 0 && s != "" && s != substr && strings.Contains(s, substr)
}
