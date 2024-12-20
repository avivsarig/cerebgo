package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

// TaskModifier defines a function that modifies a task based on the current time.
type TaskModifier func(models.Task, time.Time) models.Task

// CompletionModifier returns a TaskModifier that marks a task as completed.
// If the task is already completed, it updates the "UpdatedAt" field.
// Otherwise, it marks the task as done, sets the completion time, and updates "UpdatedAt".
//
// Parameters:
//   - completionTime: The timestamp to record as the task's completion time.
//
// Returns:
//   - TaskModifier: A function to modify a task.
func CompletionModifier(completionTime time.Time) TaskModifier {
	return func(task models.Task, now time.Time) models.Task {
		if IsCompleted(task) {
			return models.Task{
				Title:          task.Title,
				Content:        task.Content,
				IsProject:      task.IsProject,
				IsHighPriority: task.IsHighPriority,
				Done:           task.Done,
				CompletedAt:    task.CompletedAt,
				DueDate:        task.DueDate,
				DoDate:         task.DoDate,
				CreatedAt:      task.CreatedAt,
				UpdatedAt:      now,
			}
		}

		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: task.IsHighPriority,
			Done:           true,
			CompletedAt:    ptr.Some(completionTime),
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now,
		}
	}
}

// UncompleteModifier returns a TaskModifier that marks a task as not completed.
// It sets "Done" to false, clears the "CompletedAt" field, and updates "UpdatedAt".
//
// Returns:
//   - TaskModifier: A function to modify a task.
func UncompleteModifier() TaskModifier {
	return func(task models.Task, now time.Time) models.Task {
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: task.IsHighPriority,
			Done:           false,
			CompletedAt:    ptr.None[time.Time](),
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now,
		}
	}
}

// ProjectModifier returns a TaskModifier that marks a task as a project.
// It ensures "IsProject" is true and updates the "UpdatedAt" field.
//
// Parameters:
//   - task: The task to modify.
//   - now: The current timestamp.
//
// Returns:
//   - models.Task: The modified task.
func ProjectModifier(task models.Task, now time.Time) models.Task {
	return models.Task{
		Title:          task.Title,
		Content:        task.Content,
		IsProject:      true,
		IsHighPriority: task.IsHighPriority,
		CompletedAt:    task.CompletedAt,
		Done:           task.Done,
		DueDate:        task.DueDate,
		DoDate:         task.DoDate,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      now,
	}
}

// ComposeModifiers combines multiple TaskModifier functions into a single TaskModifier.
// The modifiers are applied sequentially to the task.
//
// Parameters:
//   - modifiers: A list of TaskModifier functions to apply.
//
// Returns:
//   - TaskModifier: A function that applies all the provided modifiers in order.
func ComposeModifiers(modifiers ...TaskModifier) TaskModifier {
	return func(task models.Task, now time.Time) models.Task {
		result := task
		for _, modifier := range modifiers {
			result = modifier(result, now)
		}
		return result
	}
}
