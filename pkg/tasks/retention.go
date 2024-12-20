package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/mdparser"
)

// RetentionConfig defines the retention periods for tasks based on their type.
type RetentionConfig struct {
	EmptyTaskRetention time.Duration
	ProjectRetention   time.Duration
}

// TaskRetentionResult represents the outcome of evaluating a task for retention.
type TaskRetentionResult struct {
	FilePath     string
	ShouldRetain bool
	Error        error
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

// clearCompletedTasks deletes completed task files that no longer meet retention criteria.
//
// Parameters:
//   - completedPath: The directory containing completed task files.
//   - now: The current timestamp.
//   - config: Retention rules for completed tasks.
//
// Returns:
//   - error: An error if file deletion or directory reading fails.
func clearCompletedTasks(completedPath string, now time.Time, config RetentionConfig) error {
	entries, err := os.ReadDir(completedPath)
	if err != nil {
		return fmt.Errorf("failed to read completed tasks directory: %w", err)
	}

	results := make([]TaskRetentionResult, 0, len(entries))

	for _, entry := range entries {
		result := processEntry(entry, completedPath, now, config)
		results = append(results, result)
	}

	for _, result := range results {
		if result.Error != nil {
			return result.Error
		}
		if !result.ShouldRetain {
			if err := os.Remove(result.FilePath); err != nil {
				return fmt.Errorf("failed to delete task file: %w", err)
			}
		}
	}

	return nil
}

// ClearCompletedTasks processes and clears completed tasks from the storage based on retention settings.
//
// Parameters:
//   - p: Processor containing configuration and context.
//   - now: The current timestamp.
//
// Returns:
//   - error: An error if the clearing process fails.
func ClearCompletedTasks(p Processor, now time.Time) error {
	completedPath := p.config.GetString("paths.subdirs.tasks.completed")
	config := RetentionConfig{
		EmptyTaskRetention: time.Duration(p.config.GetInt("settings.retention.empty_task")) * 24 * time.Hour,
		ProjectRetention:   time.Duration(p.config.GetInt("settings.retention.project_before_archive")) * 24 * time.Hour,
	}

	return clearCompletedTasks(completedPath, now, config)
}

// processEntry evaluates a single task file to determine if it should be retained or deleted.
//
// Parameters:
//   - entry: The file entry to process.
//   - basePath: Base directory path for the task files.
//   - now: The current timestamp.
//   - config: Retention rules for completed tasks.
//
// Returns:
//   - TaskRetentionResult: The evaluation result including file path, retention decision, and any error.
func processEntry(entry os.DirEntry, basePath string, now time.Time, config RetentionConfig) TaskRetentionResult {
	if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
		return TaskRetentionResult{ShouldRetain: true}
	}

	filePath := filepath.Join(basePath, entry.Name())
	doc, err := mdparser.ParseMarkdownDoc(filePath)
	if err != nil {
		return TaskRetentionResult{
			Error: fmt.Errorf("failed to parse task file %s: %w", entry.Name(), err),
		}
	}

	task, err := mdparser.DocumetToTask(doc)
	if err != nil {
		return TaskRetentionResult{
			Error: fmt.Errorf("failed to convert document to task %s: %w", entry.Name(), err),
		}
	}

	return TaskRetentionResult{
		FilePath:     filePath,
		ShouldRetain: ShouldRetainTask(task, now, config),
	}
}
