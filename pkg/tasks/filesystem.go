package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

// ReadTaskFile reads and parses a markdown file into a Task model
//
// Parameters:
//   - filePath: path to the markdown file
//     TODO: Check if path absolute or relative
//
// Returns:
//   - Option[Task]: Some(Task) if parsing succeeds, None if fails
//   - error: parsing or conversion errors with context
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

// readTasksFromDirectory scans a directory for markdown files and converts them to Tasks
//
// Parameters:
//   - dir: directory path containing markdown task files (.md extension)
//
// Returns:
//   - []Task: valid tasks parsed from markdown files
//   - error: reading directory or parsing task errors with context
//
// Skip files that:
//   - are directories
//   - don't have .md extension
//   - fail to parse
//   - have invalid task data
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

// DocumentToTask converts a markdown document into a Task model. It extracts task metadata
// from the document's frontmatter and content.
//
// Required frontmatter fields:
// - created_at: timestamp of task creation
// - do_date: string representing when the task should be done
//
// Optional frontmatter fields:
// - updated_at: timestamp of last update (defaults to created_at)
// - due_date: string deadline for the task
// - done: boolean indicating completion status
// - is_project: boolean marking task as a project
// - is_high_priority: boolean for priority level
//
// Returns error if required fields are missing.
func DocumentToTask(doc mdparser.MarkdownDocument) (models.Task, error) {
	fm := doc.Frontmatter

	createdAt, ok := mdparser.GetTime(fm, "created_at")
	if !ok {
		return models.Task{}, fmt.Errorf("missing required field: created_at")
	}

	doDate, ok := mdparser.GetString(fm, "do_date")
	if !ok {
		return models.Task{}, fmt.Errorf("missing required field: do_date")
	}

	task := models.Task{
		Title:     doc.Title,
		CreatedAt: createdAt,
		DoDate:    doDate,
	}

	if updatedAt, ok := mdparser.GetTime(fm, "updated_at"); ok {
		task.UpdatedAt = updatedAt
	} else {
		task.UpdatedAt = createdAt
	}

	if doc.Content != "" {
		task.Content = ptr.Some(doc.Content)
	} else {
		task.Content = ptr.None[string]()
	}

	if dueDate, ok := mdparser.GetString(fm, "due_date"); ok {
		task.DueDate = ptr.Some(dueDate)
	} else {
		task.DueDate = ptr.None[string]()
	}

	if isDone, ok := mdparser.GetBool(fm, "done"); ok {
		task.Done = isDone
	}

	if isProject, ok := mdparser.GetBool(fm, "is_project"); ok {
		task.IsProject = isProject
	}

	if isHighPriority, ok := mdparser.GetBool(fm, "is_high_priority"); ok {
		task.IsHighPriority = isHighPriority
	}

	return task, nil
}
