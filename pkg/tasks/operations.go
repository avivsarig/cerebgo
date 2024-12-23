package tasks

import (
	"fmt"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/files"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

// TaskModifier defines a function that modifies a task based on the current time.
type TaskModifier func(models.Task, time.Time) (models.Task, error)

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
	return func(task models.Task, now time.Time) (models.Task, error) {
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
				UpdatedAt:      task.CompletedAt.Value(), // Don't update timestamp
			}, nil
		}

		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: task.IsHighPriority,
			Done:           true,                     // mark as done
			CompletedAt:    ptr.Some(completionTime), // set completion time
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now, // update timestamp
		}, nil
	}
}

// UncompleteModifier returns a TaskModifier that marks a task as not completed.
// It sets "Done" to false, clears the "CompletedAt" field, and updates "UpdatedAt".
//
// Returns:
//   - TaskModifier: A function to modify a task.
func UncompleteModifier() TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: task.IsHighPriority,
			Done:           false,                 // mark as not done
			CompletedAt:    ptr.None[time.Time](), // clear completion time
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now, // update timestamp
		}, nil
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
func ProjectModifier(now time.Time) TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      true, // mark as project
			IsHighPriority: task.IsHighPriority,
			CompletedAt:    task.CompletedAt,
			Done:           task.Done,
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now, // update timestamp
		}, nil
	}
}

func UnprojectModifier(now time.Time) TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      false, // mark as not project
			IsHighPriority: task.IsHighPriority,
			CompletedAt:    task.CompletedAt,
			Done:           task.Done,
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now, // update timestamp
		}, nil
	}
}

// DeleteModifier returns a TaskModifier that deletes a task.
//
// Parameters:
// - path: The path to the task file.
//
// Returns:
// - TaskModifier: A function to delete a task that returns an empty task.
func DeleteModifier(path string) TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		if err := DeleteTaskFile(task, path); err != nil {
			return models.Task{}, fmt.Errorf("failed to delete task file: %w", err)
		}
		return models.Task{}, nil
	}
}

// DeactivateModifier returns a TaskModifier that moves a task from active to completed.
//
// Returns:
//   - TaskModifier: A function to deactivate a task.
//     The function moves the task file from the active directory to the completed directory.
func DeactivateModifier() TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		err := files.MoveFile(
			files.FilePath{
				Dir:  configuration.GetString("paths.subdir.tasks.active"),
				Name: task.Title + ".md",
			},
			files.FilePath{
				Dir:  configuration.GetString("paths.subdir.tasks.completed"),
				Name: task.Title + ".md",
			},
		)

		if err != nil {
			return models.Task{}, fmt.Errorf("failed to move task file: %w", err)
		}

		return models.Task{}, nil
	}
}

func ReactivateModifier() TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		err := files.MoveFile(
			files.FilePath{
				Dir:  configuration.GetString("paths.subdir.tasks.completed"),
				Name: task.Title + ".md",
			},
			files.FilePath{
				Dir:  configuration.GetString("paths.subdir.tasks.active"),
				Name: task.Title + ".md",
			},
		)

		if err != nil {
			return models.Task{}, fmt.Errorf("failed to move task file: %w", err)
		}

		return models.Task{}, nil
	}
}

// DoDateTodayModifier returns a TaskModifier that sets a task's DoDate to today's date.
//
// The modifier preserves all other task fields while updating:
// - DoDate: Set to today's date in YYYY-MM-DD format
// - UpdatedAt: Set to provided timestamp
//
// Returns:
//
//	TaskModifier function that takes (models.Task, time.Time) and returns modified models.Task
func DoDateTodayModifier() TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		today := time.Now().Format("2006-01-02")
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: task.IsHighPriority,
			CompletedAt:    task.CompletedAt,
			Done:           task.Done,
			DueDate:        task.DueDate,
			DoDate:         today,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now,
		}, nil
	}
}

// HighPriorityModifier returns a TaskModifier that marks a task as high priority.
//
// Returns:
//   - TaskModifier: A function to mark a task as high priority.
//     The function sets the "IsHighPriority" field to true and updates the "UpdatedAt" field.
func HighPriorityModifier() TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		return models.Task{
			Title:          task.Title,
			Content:        task.Content,
			IsProject:      task.IsProject,
			IsHighPriority: true,
			CompletedAt:    task.CompletedAt,
			Done:           task.Done,
			DueDate:        task.DueDate,
			DoDate:         task.DoDate,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      now,
		}, nil
	}
}

// TODO: ArchiveModifier - convert to archive type (not implemented), move to archive directory

// ComposeModifiers combines multiple TaskModifier functions into a single TaskModifier.
// The modifiers are applied sequentially to the task.
//
// Parameters:
//   - modifiers: A list of TaskModifier functions to apply.
//
// Returns:
//   - TaskModifier: A function that applies all the provided modifiers in order.
func ComposeModifiers(modifiers ...TaskModifier) TaskModifier {
	return func(task models.Task, now time.Time) (models.Task, error) {
		result := task
		for _, modifier := range modifiers {
			var err error
			result, err = modifier(result, now)
			if err != nil {
				return models.Task{}, fmt.Errorf("modifier failed: %w", err)
			}
		}
		return result, nil
	}
}
