package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/spf13/viper"
)

// CreateTestDirectory creates a temporary directory for testing.
// It registers a cleanup function to remove the directory after the test.
func CreateTestDirectory(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("failed to clean up test directory: %v", err)
		}
	})

	return dir
}

// CreateTestFile creates a file with given content in the specified directory.
func CreateTestFile(t *testing.T, dir string, filename string, content string) error {
	t.Helper()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %v", err)
	}

	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// CreateTestTaskFile creates a markdown file for a task.
func CreateTestTaskFile(t *testing.T, basePath string, filename string, task models.Task) error {
	t.Helper()
	content := generateTaskMarkdown(task)
	return CreateTestFile(t, basePath, filename, content)
}

// AssertFileExists checks if a file exists and fails the test if it doesn't.
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s, but it doesn't", path)
	}
}

// AssertFileNotExists checks if a file doesn't exist and fails the test if it does.
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to not exist at %s, but it does", path)
	}
}

// AssertFileContent checks if a file's content matches the expected content.
func AssertFileContent(t *testing.T, path string, expected string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("failed to read file at %s: %v", path, err)
		return
	}

	if string(content) != expected {
		t.Errorf("file content mismatch at %s\ngot: %s\nwant: %s",
			path, content, expected)
	}
}

// MoveTestFile moves a file in a test environment.
// It fails the test if the operation fails.
func MoveTestFile(t *testing.T, sourcePath, destPath string) {
	t.Helper()

	if err := os.Rename(sourcePath, destPath); err != nil {
		t.Errorf("failed to move file from %s to %s: %v",
			sourcePath, destPath, err)
	}
}

// DeleteTestFile deletes a file in a test environment.
// It fails the test if the operation fails.
func DeleteTestFile(t *testing.T, path string) {
	t.Helper()

	if err := os.Remove(path); err != nil {
		t.Errorf("failed to delete file at %s: %v", path, err)
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

// setupConfigDir creates a temporary directory and config file for testing.
// It returns the directory path. The directory is automatically cleaned up after the test.
func SetupConfigDir(t *testing.T, content string) string {
	t.Helper()
	dir := CreateTestDirectory(t)
	err := CreateTestFile(t, dir, "config.yaml", content)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

// setConfigPath sets the CONFIG_PATH environment variable for testing and
// ensures it's restored to its original value after the test completes.
func SetConfigPath(t *testing.T, path string) {
	t.Helper()
	oldEnv := os.Getenv("CONFIG_PATH")
	os.Setenv("CONFIG_PATH", path)
	t.Cleanup(func() {
		os.Setenv("CONFIG_PATH", oldEnv)
	})
}

// validateConfig checks if the loaded configuration matches expected values.
// It verifies both top-level and nested configuration settings.
func ValidateConfig(t *testing.T, v *viper.Viper) {
	t.Helper()
	if v.GetString("key") != "value" {
		t.Errorf("key = %v, want %v", v.GetString("key"), "value")
	}
	if v.GetInt("nested.setting") != 42 {
		t.Errorf("nested.setting = %v, want %v", v.GetInt("nested.setting"), 42)
	}
}
