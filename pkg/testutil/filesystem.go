package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
)

func CreateTestTaskFile(t *testing.T, basePath string, filename string, task models.Task) error {
	t.Helper()

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %v", err)
	}

	content := generateTaskMarkdown(task)

	fullPath := filepath.Join(basePath, filename)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write task file: %v", err)
	}

	return nil
}

func AssertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s, but it doesn't", path)
	}
}

func generateTaskMarkdown(task models.Task) string {
	// Start with frontmatter
	content := fmt.Sprintf(`---
created_at: %s
updated_at: %s
do_date: %s
`,
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
		task.DoDate,
	)

	if task.Content.IsValid() {
		content += fmt.Sprintf("content: %s\n", task.Content.Value())
	}
	if task.DueDate.IsValid() {
		content += fmt.Sprintf("due_date: %s\n", task.DueDate.Value())
	}
	if task.Done {
		content += "done: true\n"
		if task.CompletedAt.IsValid() {
			content += fmt.Sprintf("completed_at: %s\n",
				task.CompletedAt.Value().Format(time.RFC3339))
		}
	}
	if task.IsProject {
		content += "is_project: true\n"
	}
	if task.IsHighPriority {
		content += "is_high_priority: true\n"
	}

	content += "---\n\n"
	if task.Content.IsValid() {
		content += task.Content.Value()
	}

	return content
}
