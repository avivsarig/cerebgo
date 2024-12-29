package tasks_test

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/tasks"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// TestShouldRetainTask uses table-driven testing to verify the retention logic.
func TestShouldRetainTask(t *testing.T) {
	// Define a fixed "now" time for consistent testing
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Define retention config for tests
	config := tasks.RetentionConfig{
		EmptyTaskRetention: 7 * 24 * time.Hour,  // 7 days
		ProjectRetention:   30 * 24 * time.Hour, // 30 days
	}

	// Helper function to create a completed task
	createCompletedTask := func(isProject bool, completedAt time.Time) models.Task {
		return models.Task{
			Done:        true,
			IsProject:   isProject,
			CompletedAt: ptr.Some(completedAt),
		}
	}

	tests := []struct {
		name    string
		task    models.Task
		want    bool
		message string
	}{
		{
			name: "uncompleted_task_should_be_retained",
			task: models.Task{
				Done: false,
			},
			want:    true,
			message: "uncompleted tasks should always be retained",
		},
		{
			name:    "recently_completed_task_should_be_retained",
			task:    createCompletedTask(false, now.Add(-6*24*time.Hour)),
			want:    true,
			message: "tasks completed within retention period should be retained",
		},
		{
			name:    "old_completed_task_should_not_be_retained",
			task:    createCompletedTask(false, now.Add(-8*24*time.Hour)),
			want:    false,
			message: "tasks completed beyond retention period should not be retained",
		},
		{
			name:    "recently_completed_project_should_be_retained",
			task:    createCompletedTask(true, now.Add(-29*24*time.Hour)),
			want:    true,
			message: "projects completed within retention period should be retained",
		},
		{
			name:    "old_completed_project_should_not_be_retained",
			task:    createCompletedTask(true, now.Add(-31*24*time.Hour)),
			want:    false,
			message: "projects completed beyond retention period should not be retained",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tasks.ShouldRetainTask(tt.task, now, config)
			result := testutil.ValidateEqual(
				"retention decision",
				got,
				tt.want,
			)

			if !result.IsValid {
				t.Errorf("%s: %s", tt.message, result.ToString())
			}
		})
	}
}
