package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

type TaskModifier func(models.Task, time.Time) models.Task

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

func ComposeModifiers(modifiers ...TaskModifier) TaskModifier {
	return func(task models.Task, now time.Time) models.Task {
		result := task
		for _, modifier := range modifiers {
			result = modifier(result, now)
		}
		return result
	}
}
