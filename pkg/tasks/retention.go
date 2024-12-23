package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
)

// RetentionConfig defines the retention periods for tasks based on their type.
type RetentionConfig struct {
	EmptyTaskRetention time.Duration
	ProjectRetention   time.Duration
}

// ShouldRetainTask checks whether a completed task should be retained based on its type and age.
//
// Parameters:
//   - task: The task to evaluate.
//   - now: The current timestamp.
//   - config: Retention rules for completed tasks.
//
// Returns:
//   - bool: True if the task should be retained, false otherwise.
func ShouldRetainTask(task models.Task, now time.Time, config RetentionConfig) bool {
	if !IsCompleted(task) {
		return true
	}

	completedAge := now.Sub(task.CompletedAt.Value())
	if task.IsProject {
		return completedAge <= config.ProjectRetention
	}
	return completedAge <= config.EmptyTaskRetention
}
