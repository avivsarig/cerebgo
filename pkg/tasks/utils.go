package tasks

import "github.com/avivSarig/cerebgo/internal/models"

// IsCompleted checks if task if fully completed
//
// Parameters:
//   - task: Task to check
//
// Returns:
//   - bool: true if task is Done and has valid completed_at timestamp
func IsCompleted(t models.Task) bool {
	return t.Done && t.CompletedAt.IsValid()
}
