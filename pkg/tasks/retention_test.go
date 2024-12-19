package tasks

import (
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// To avoid filesystem operations in tests.
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

type Config interface {
	GetString(key string) string
	GetInt(key string) int
}

type mockConfig struct {
	paths    map[string]string
	settings map[string]int
}

func (m *mockConfig) GetString(key string) string {
	return m.paths[key]
}

func (m *mockConfig) GetInt(key string) int {
	return m.settings[key]
}

type mockProcessor struct {
	config Config
}

func newMockProcessor(config Config) *mockProcessor {
	return &mockProcessor{
		config: config,
	}
}

func (p *mockProcessor) ClearCompletedTasks(now time.Time) error {
	completedPath := p.config.GetString("paths.subdirs.tasks.completed")
	config := RetentionConfig{
		EmptyTaskRetention: time.Duration(p.config.GetInt("settings.retention.empty_task")) * 24 * time.Hour,
		ProjectRetention:   time.Duration(p.config.GetInt("settings.retention.project_before_archive")) * 24 * time.Hour,
	}

	return clearCompletedTasks(completedPath, now, config)
}

func TestShouldRetainTask(t *testing.T) {
	now := time.Now()
	config := RetentionConfig{
		EmptyTaskRetention: 30 * 24 * time.Hour, // 30 days
		ProjectRetention:   7 * 24 * time.Hour,  // 7 days
	}

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "incomplete task should be retained",
			task: models.Task{
				Done:        false,
				CompletedAt: ptr.None[time.Time](),
			},
			want: true,
		},
		{
			name: "empty task older than retention should not be retained",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.Some(now.Add(-31 * 24 * time.Hour)),
				IsProject:   false,
			},
			want: false,
		},
		{
			name: "empty task within retention should be retained",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.Some(now.Add(-29 * 24 * time.Hour)),
				IsProject:   false,
			},
			want: true,
		},
		{
			name: "project older than retention should not be retained",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.Some(now.Add(-8 * 24 * time.Hour)),
				IsProject:   true,
			},
			want: false,
		},
		{
			name: "project within retention should be retained",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.Some(now.Add(-6 * 24 * time.Hour)),
				IsProject:   true,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldRetainTask(tt.task, now, config)
			results := []testutil.ValidationResult{
				testutil.ValidateEqual("retention decision", got, tt.want),
			}
			testutil.ReportResults(t, results)
		})
	}
}

func TestProcessEntry(t *testing.T) {
	now := time.Now()
	config := RetentionConfig{
		EmptyTaskRetention: 30 * 24 * time.Hour,
		ProjectRetention:   7 * 24 * time.Hour,
	}

	testDir := t.TempDir()

	tests := []struct {
		name           string
		entry          fs.DirEntry
		expectedResult TaskRetentionResult
	}{
		{
			name:  "directory entry should be retained",
			entry: mockDirEntry{name: "test-dir", isDir: true},
			expectedResult: TaskRetentionResult{
				ShouldRetain: true,
			},
		},
		{
			name:  "non-markdown file should be retained",
			entry: mockDirEntry{name: "test.txt", isDir: false},
			expectedResult: TaskRetentionResult{
				ShouldRetain: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processEntry(tt.entry, testDir, now, config)

			results := []testutil.ValidationResult{
				testutil.ValidateEqual("should retain", got.ShouldRetain, tt.expectedResult.ShouldRetain),
			}
			testutil.ReportResults(t, results)
		})
	}
}

func TestClearCompletedTasks(t *testing.T) {
	tests := []struct {
		name        string
		config      mockConfig
		setupFiles  func(dir string) error
		wantErr     bool
		errContains string
	}{
		{
			name: "invalid directory path should return error",
			config: mockConfig{
				paths: map[string]string{
					"paths.subdirs.tasks.completed": "/nonexistent/path",
				},
				settings: map[string]int{
					"settings.retention.empty_task":             30,
					"settings.retention.project_before_archive": 7,
				},
			},
			setupFiles:  func(dir string) error { return nil },
			wantErr:     true,
			errContains: "failed to read completed tasks directory",
		},
		{
			name: "empty directory should succeed",
			config: mockConfig{
				paths: map[string]string{
					"paths.subdirs.tasks.completed": "", // Will be set to temp dir
				},
				settings: map[string]int{
					"settings.retention.empty_task":             30,
					"settings.retention.project_before_archive": 7,
				},
			},
			setupFiles: func(dir string) error { return nil },
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()

			if tt.config.paths["paths.subdirs.tasks.completed"] == "" {
				tt.config.paths["paths.subdirs.tasks.completed"] = testDir
			}

			if err := tt.setupFiles(testDir); err != nil {
				t.Fatal(err)
			}

			p := newMockProcessor(&tt.config)
			err := p.ClearCompletedTasks(time.Now())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
