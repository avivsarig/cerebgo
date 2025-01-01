package tasks

import (
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
)

// PlanCompletedTaskActions plans the actions to take on a completed task.
//
// Parameters:
//   - task: The task to process.
//   - now: The current timestamp.
//   - config: The retention configuration.
//
// Returns:
//   - []TaskAction: The actions to take on the task.
//   - error: An error if the actions cannot be planned.
func PlanCompletedTaskActions(task models.Task, now time.Time, config RetentionConfig) ([]TaskModifier, error) {
	dir := configuration.GetString("paths.subdirs.tasks.completed")
	actions := make([]TaskModifier, 0)

	if !IsCompleted(task) {
		// is partial completed?
		if task.Done {
			actions = append(actions, CompletionModifier(now))
		}

		// is completion but not done?
		if task.CompletedAt.IsValid() && !task.Done {
			actions = append(actions, UncompleteModifier())
			actions = append(actions, ReactivateModifier())
		}
	}

	if !ShouldRetainTask(task, now, config) {
		if task.IsProject {
			actions = append(actions, ArchiveModifier())
		}
		actions = append(actions, DeleteModifier(dir))
	}

	return actions, nil
}

// PlanActiveTaskActions plans the actions to take on an active task.
//
// Parameters:
//   - task: The task to process.
//   - now: The current timestamp.
//   - config: The retention configuration.
//
// Returns:
//   - []TaskAction: The actions to take on the task.
//   - error: An error if the actions cannot be planned.
func PlanActiveTaskActions(task models.Task, now time.Time) ([]TaskModifier, error) {
	actions := make([]TaskModifier, 0)

	if task.Content.IsValid() && !task.IsProject {
		actions = append(actions, ProjectModifier(now))
	}
	if !task.Content.IsValid() && task.IsProject {
		actions = append(actions, UnprojectModifier(now))
	}

	if IsValidDoDate(task, now) {
		actions = append(actions, DoDateTodayModifier())
	}

	// FUTURE: Handle high priority according to DueDate

	if task.Done {
		actions = append(actions, CompletionModifier(now))
		actions = append(actions, DeactivateModifier())
	}

	return actions, nil
}
