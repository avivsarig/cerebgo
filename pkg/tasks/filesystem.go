package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func ReadTaskFile(filePath string) (ptr.Option[models.Task], error) {
	doc, err := mdparser.ParseMarkdownDoc(filePath)
	if err != nil {
		return ptr.None[models.Task](), fmt.Errorf("failed to parse markdown from %s: %w", filePath, err)
	}

	task, err := DocumentToTask(doc)
	if err != nil {
		return ptr.None[models.Task](), fmt.Errorf("failed to convert document to task from %s: %w", filePath, err)
	}

	return ptr.Some(task), nil
}

func readTasksFromDirectory(dir string) ([]models.Task, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var tasks []models.Task
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		taskResult, err := ReadTaskFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read task from %s: %w", filePath, err)
		}

		if taskResult.IsValid() {
			tasks = append(tasks, taskResult.Value())
		}
	}

	return tasks, nil
}

// func writeTaskToFile(task models.Task, path string) error {
// 	// TODO: Implement task to markdown conversion.
// 	// For now we'll leave this as a stub until we implement the conversion logic.
// 	return nil
// }

func DocumentToTask(doc mdparser.MarkdownDocument) (models.Task, error) {
	fp := mdparser.NewFrontmatterProcessor(doc.Frontmatter)

	createdAt, ok := fp.GetTime("created_at")
	if !ok {
		return models.Task{}, fmt.Errorf("missing required field: created_at")
	}

	doDate, ok := fp.GetString("do_date")
	if !ok {
		return models.Task{}, fmt.Errorf("missing required field: do_date")
	}

	task := models.Task{
		Title:     doc.Title,
		CreatedAt: createdAt,
		DoDate:    doDate,
	}

	if updatedAt, ok := fp.GetTime("updated_at"); ok {
		task.UpdatedAt = updatedAt
	} else {
		task.UpdatedAt = createdAt
	}

	if doc.Content != "" {
		task.Content = ptr.Some(doc.Content)
	} else {
		task.Content = ptr.None[string]()
	}

	if dueDate, ok := fp.GetString("due_date"); ok {
		task.DueDate = ptr.Some(dueDate)
	} else {
		task.DueDate = ptr.None[string]()
	}

	if isDone, ok := fp.GetBool("done"); ok {
		task.Done = isDone
	}

	if isProject, ok := fp.GetBool("is_project"); ok {
		task.IsProject = isProject
	}

	if isHighPriority, ok := fp.GetBool("is_high_priority"); ok {
		task.IsHighPriority = isHighPriority
	}

	if task.Done {
		if completedAt, ok := fp.GetTime("completed_at"); ok {
			task.CompletedAt = ptr.Some(completedAt)
		} else {
			return models.Task{}, fmt.Errorf("done task missing completed_at date")
		}
	} else {
		task.CompletedAt = ptr.None[time.Time]()
	}

	return task, nil
}
