package mdparser_test

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestDocumentToTask(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	baseTimeStr := baseTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		doc     mdparser.MarkdownDocument
		want    models.Task
		wantErr bool
	}{
		{
			name: "valid document with all fields",
			doc: mdparser.MarkdownDocument{
				Title: "Test Task",
				Frontmatter: mdparser.Frontmatter{
					"created_at": baseTimeStr,
					"updated_at": baseTimeStr,
					"due_date":   baseTimeStr,
					"priority":   "high",
				},
				Content: "Test content",
			},
			want: models.Task{
				Title:          "Test Task",
				Content:        ptr.Some("Test content"),
				IsProject:      true,
				IsHighPriority: true,
				Done:           false,
				CreatedAt:      baseTime,
				UpdatedAt:      baseTime,
				DueDate:        ptr.Some(baseTimeStr),
				CompletedAt:    ptr.None[time.Time](),
			},
		},
		{
			name: "minimal valid document",
			doc: mdparser.MarkdownDocument{
				Title: "Test Task",
				Frontmatter: mdparser.Frontmatter{
					"created_at": baseTimeStr,
					"updated_at": baseTimeStr,
				},
			},
			want: models.Task{
				Title:          "Test Task",
				Content:        ptr.Some(""),
				IsProject:      false,
				IsHighPriority: false,
				Done:           false,
				CreatedAt:      baseTime,
				UpdatedAt:      baseTime,
				DueDate:        ptr.None[string](),
				CompletedAt:    ptr.None[time.Time](),
			},
		},
		{
			name: "missing created_at",
			doc: mdparser.MarkdownDocument{
				Title: "Test Task",
				Frontmatter: mdparser.Frontmatter{
					"updated_at": baseTimeStr,
				},
			},
			wantErr: true,
		},
		{
			name: "missing title",
			doc: mdparser.MarkdownDocument{
				Frontmatter: mdparser.Frontmatter{
					"created_at": baseTimeStr,
					"updated_at": baseTimeStr,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			doc: mdparser.MarkdownDocument{
				Title: "Test Task",
				Frontmatter: mdparser.Frontmatter{
					"created_at": "2024-01-01", // Not RFC3339
					"updated_at": baseTimeStr,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid priority value",
			doc: mdparser.MarkdownDocument{
				Title: "Test Task",
				Frontmatter: mdparser.Frontmatter{
					"created_at": baseTimeStr,
					"updated_at": baseTimeStr,
					"priority":   "SUPER HIGH", // Invalid priority
				},
			},
			want: models.Task{
				Title:          "Test Task",
				Content:        ptr.Some(""),
				IsProject:      false,
				IsHighPriority: false, // Should default to false for invalid priority
				Done:           false,
				CreatedAt:      baseTime,
				UpdatedAt:      baseTime,
				DueDate:        ptr.None[string](),
				CompletedAt:    ptr.None[time.Time](),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mdparser.DocumentToTask(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("DocumentToTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				testutil.AssertTaskEqual(t, got, tt.want)
			}
		})
	}
}

func TestParseMarkdownDocFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		want     mdparser.MarkdownDocument
		wantErr  bool
	}{
		{
			name: "valid frontmatter and content",
			content: `---
key: value
date: 2024-01-01
---

# Content here`,
			filename: "test.md",
			want: mdparser.MarkdownDocument{
				Title: "test",
				Frontmatter: mdparser.Frontmatter{
					"key":  "value",
					"date": "2024-01-01",
				},
				Content: "# Content here",
			},
		},
		{
			name:     "empty file",
			content:  "",
			filename: "empty.md",
			want: mdparser.MarkdownDocument{
				Title: "empty",
			},
		},
		{
			name: "invalid frontmatter",
			content: `---
invalid: [yaml
---`,
			filename: "invalid.md",
			wantErr:  true,
		},
		{
			name: "multiple frontmatter blocks",
			content: `---
key: value
---
content
---
more: stuff
---`,
			filename: "multiple.md",
			wantErr:  true,
		},
		{
			name: "content before frontmatter",
			content: `Some content
---
key: value
---`,
			filename: "invalid.md",
			wantErr:  true,
		},
		{
			name: "whitespace before frontmatter",
			content: ` ---
key: value
---`,
			filename: "invalid.md",
			wantErr:  true,
		},
		{
			name: "unclosed frontmatter",
			content: `---
key: value
content`,
			filename: "unclosed.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile := testutil.CreateTestDirectory(t)
			filePath := tmpFile + "/" + tt.filename
			err := testutil.CreateTestFile(t, tmpFile, tt.filename, tt.content)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			got, err := mdparser.ParseMarkdownDoc(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Title != tt.want.Title {
					t.Errorf("Title = %v, want %v", got.Title, tt.want.Title)
				}
				if got.Content != tt.want.Content {
					t.Errorf("Content = %v, want %v", got.Content, tt.want.Content)
				}
				// Compare frontmatter
				if len(got.Frontmatter) != len(tt.want.Frontmatter) {
					t.Errorf("Frontmatter length = %v, want %v", len(got.Frontmatter), len(tt.want.Frontmatter))
				}
				for k, v := range tt.want.Frontmatter {
					if got.Frontmatter[k] != v {
						t.Errorf("Frontmatter[%v] = %v, want %v", k, got.Frontmatter[k], v)
					}
				}
			}
		})
	}
}
