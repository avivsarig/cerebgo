package tasks

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestHandleTaskAction(t *testing.T) {
	// Setup test time
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	testBasePath := t.TempDir()

	// Create a helper function that generates a task for testing
	makeTestTask := func() models.Task {
		return models.Task{
			Title:     "test-task",
			Content:   ptr.Some("test content"),
			CreatedAt: now,
			UpdatedAt: now,
			DoDate:    "2024-01-01",
		}
	}

	tests := []struct {
		name  string
		input struct {
			action   TaskAction
			basePath string
			now      time.Time
		}
		setup    func(t *testing.T, basePath string, task models.Task) error
		validate func(t *testing.T, basePath string) error
		wantErr  bool
		errMsg   string
	}{
		{
			name: "retain task - success",
			input: struct {
				action   TaskAction
				basePath string
				now      time.Time
			}{
				action: TaskAction{
					Task:   makeTestTask(),
					Action: ActionRetain,
				},
				basePath: testBasePath,
				now:      now,
			},
			setup: func(t *testing.T, basePath string, task models.Task) error {
				return testutil.CreateTestTaskFile(t, basePath, "test-task.md", task)
			},
			validate: func(t *testing.T, basePath string) error {
				testutil.AssertFileExists(t, filepath.Join(basePath, "test-task.md"))
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We now pass the task to setup
			if tt.setup != nil {
				if err := tt.setup(t, tt.input.basePath, tt.input.action.Task); err != nil {
					t.Fatalf("test setup failed: %v", err)
				}
			}

			err := handleTaskAction(tt.input.action, tt.input.basePath, tt.input.now)

			if (err != nil) != tt.wantErr {
				t.Errorf("handleTaskAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("handleTaskAction() error message = %v, want %v", err, tt.errMsg)
				}
			}

			if tt.validate != nil {
				if err := tt.validate(t, tt.input.basePath); err != nil {
					t.Errorf("validation failed: %v", err)
				}
			}
		})
	}
}
