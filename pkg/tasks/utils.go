package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
)

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

// IsValidDoDate checks if the task's DoDate is valid - has valid format and is not in the past
//
// Parameters:
//   - task: Task to check
//
// Returns:
//   - bool: true if task's DoDate is valid
func IsValidDoDate(t models.Task, now time.Time) bool {
	doDateStr, err := time.Parse("2006-01-02", t.DoDate)
	if err != nil {
		return false
	}

	nowDate := now.Truncate(24 * time.Hour)
	doDate := doDateStr.Truncate(24 * time.Hour)

	return !nowDate.After(doDate)
}

// IsValidDueDate checks if the task's DueDate is valid - has valid format and is not in the past
//
// Parameters:
//   - task: Task to check
//
// Returns:
//   - bool: true if task's DoDate is valid
func IsValidDueDate(t models.Task, now time.Time) bool {
	if !t.DueDate.IsValid() {
		return false
	}

	dueDateStr, err := time.Parse("2006-01-02", t.DueDate.Value())
	if err != nil {
		return false
	}

	nowDate := now.Truncate(24 * time.Hour)
	dueDate := dueDateStr.Truncate(24 * time.Hour)

	return !nowDate.After(dueDate)
}
