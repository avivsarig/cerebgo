package records_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/records"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// TestWriteRecordToFile_Creation tests the basic file creation functionality.
func TestWriteRecordToFile_Creation(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		record   models.Record
		setup    func(t *testing.T, dir string)                       // For complex test setup
		validate func(t *testing.T, dir string, record models.Record) // For complex validations
		wantErr  bool
	}{
		{
			name: "basic record with minimum fields",
			record: models.Record{
				Title:     "test-record",
				CreatedAt: baseTime,
				UpdatedAt: baseTime,
				Tags:      []string{"tag1"},
			},
		},
		// FUTURE:
		// This test is commented out because the current implementation does not sanitize file names.
		// {
		// 	name: "record with special characters in title",
		// 	record: models.Record{
		// 		Title:     "test/record:with?special*chars",
		// 		CreatedAt: baseTime,
		// 		UpdatedAt: baseTime,
		// 		Tags:      []string{"tag1"},
		// 	},
		// 	validate: func(t *testing.T, dir string, record models.Record) {
		// 		// Verify file name is sanitized
		// 		sanitized := records.sanitizeFileName(record.Title)
		// 		filePath := filepath.Join(dir, sanitized+".md")
		// 		testutil.AssertFileExists(t, filePath)
		// 	},
		// },
		{
			name: "record with all optional fields",
			record: models.Record{
				Title:      "full-record",
				CreatedAt:  baseTime,
				UpdatedAt:  baseTime,
				Tags:       []string{"tag1", "tag2", "tag3"},
				Content:    ptr.Some("test content\nwith multiple\nlines"),
				URL:        ptr.Some("https://example.com"),
				ArchivedAt: ptr.Some(baseTime),
			},
		},
		{
			name: "record with empty optional fields",
			record: models.Record{
				Title:      "empty-fields",
				CreatedAt:  baseTime,
				UpdatedAt:  baseTime,
				Tags:       []string{}, // Empty tags
				Content:    ptr.None[string](),
				URL:        ptr.None[string](),
				ArchivedAt: ptr.None[time.Time](),
			},
		},
		{
			name: "existing file should be overwritten",
			record: models.Record{
				Title:     "existing-record",
				CreatedAt: baseTime,
				UpdatedAt: baseTime,
				Tags:      []string{"new-tag"},
			},
			setup: func(t *testing.T, dir string) {
				err := testutil.CreateTestFile(t, dir, "existing-record.md", "old content")
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			},
			validate: func(t *testing.T, dir string, record models.Record) {
				content, err := os.ReadFile(filepath.Join(dir, "existing-record.md"))
				if err != nil {
					t.Fatalf("Failed to read file: %v", err)
				}
				if strings.Contains(string(content), "old content") {
					t.Error("File should have been overwritten")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := testutil.CreateTestDirectory(t)

			if tt.setup != nil {
				tt.setup(t, dir)
			}

			err := records.WriteRecordToFile(tt.record, dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRecordToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Default validation if no custom validation provided
			if tt.validate == nil {
				filePath := filepath.Join(dir, tt.record.Title+".md")
				testutil.AssertFileExists(t, filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read created file: %v", err)
				}

				// Basic content validations
				if !strings.Contains(string(content), "created_at: "+tt.record.CreatedAt.Format(time.RFC3339)) {
					t.Error("Created timestamp not found in file content")
				}

				for _, tag := range tt.record.Tags {
					if !strings.Contains(string(content), "- "+tag) {
						t.Errorf("Tag %q not found in file content", tag)
					}
				}

				if tt.record.Content.IsValid() && !strings.Contains(string(content), tt.record.Content.Value()) {
					t.Error("Content not found in file")
				}
			} else {
				tt.validate(t, dir, tt.record)
			}
		})
	}
}

// TestWriteRecordToFile_Errors tests error conditions.
func TestWriteRecordToFile_Errors(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		record  models.Record
		path    string
		wantErr bool
	}{
		{
			name: "invalid path",
			record: models.Record{
				Title:     "test",
				CreatedAt: baseTime,
				UpdatedAt: baseTime,
				Tags:      []string{"tag1"},
			},
			path:    "/nonexistent/path",
			wantErr: true,
		},
		{
			name: "empty title",
			record: models.Record{
				Title:     "",
				CreatedAt: baseTime,
				UpdatedAt: baseTime,
				Tags:      []string{"tag1"},
			},
			wantErr: true,
		},
		{
			name: "read-only directory",
			record: models.Record{
				Title:     "test",
				CreatedAt: baseTime,
				UpdatedAt: baseTime,
				Tags:      []string{"tag1"},
			},
			// Path will be made read-only in the test
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dir string
			if tt.path != "" {
				dir = tt.path
			} else {
				dir = testutil.CreateTestDirectory(t)

				if tt.name == "read-only directory" {
					// Make directory read-only
					err := os.Chmod(dir, 0444)
					if err != nil {
						t.Fatalf("Failed to make directory read-only: %v", err)
					}
					// Restore permissions after test
					defer os.Chmod(dir, 0755)
				}
			}

			err := records.WriteRecordToFile(tt.record, dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRecordToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
