package tasks

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func TestCompletionModifier(t *testing.T) {
	now := time.Now()
	completionTime := now.Add(-1 * time.Hour)

	task := models.Task{
		Title:          "Test Task",
		Content:        ptr.Some("content"),
		IsProject:      false,
		IsHighPriority: true,
		CompletedAt:    ptr.None[time.Time](),
		CreatedAt:      now.Add(-2 * time.Hour),
	}

	modified := CompletionModifier(completionTime)(task, now)

	if modified.Title != task.Title {
		t.Errorf("Title changed: got %v, want %v", modified.Title, task.Title)
	}

	if !modified.CompletedAt.IsValid() || !modified.CompletedAt.Value().Equal(completionTime) {
		t.Error("CompletedAt not set correctly")
	}

	if !modified.UpdatedAt.Equal(now) {
		t.Error("UpdatedAt not set correctly")
	}
}

func TestProjectModifier(t *testing.T) {
	now := time.Now()
	task := models.Task{
		Title:     "Test Task",
		IsProject: false,
	}

	modified := ProjectModifier(task, now)

	if !modified.IsProject {
		t.Error("Task not marked as project")
	}

	if !modified.UpdatedAt.Equal(now) {
		t.Error("UpdatedAt not set correctly")
	}
}

func TestComposeModifiers(t *testing.T) {
	now := time.Now()
	completionTime := now.Add(-1 * time.Hour)

	task := models.Task{
		Title:     "Test Task",
		IsProject: false,
	}

	composed := ComposeModifiers(
		ProjectModifier,
		CompletionModifier(completionTime),
	)

	modified := composed(task, now)

	if !modified.IsProject {
		t.Error("Project modifier not applied")
	}

	if !modified.CompletedAt.IsValid() || !modified.CompletedAt.Value().Equal(completionTime) {
		t.Error("Completion modifier not applied")
	}
}
