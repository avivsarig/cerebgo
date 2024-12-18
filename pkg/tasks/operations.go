package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

type TaskModifier func(models.Task, time.Time) models.Task

func CompletionModifier(completionTime time.Time) TaskModifier {
	return func(task models.Task, now time.Time) models.Task {
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

func ProjectModifier(task models.Task, now time.Time) models.Task {
	return models.Task{
		Title:          task.Title,
		Content:        task.Content,
		IsProject:      true,
		IsHighPriority: task.IsHighPriority,
		CompletedAt:    task.CompletedAt,
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
