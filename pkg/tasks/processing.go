package tasks

import (
	"fmt"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

// getTaskUpdates determines the necessary updates for a given task based on its state.
// If the task has valid content and is not a project, it suggests a project update.
// If the task is marked as done, it schedules a completion update.
//
// Parameters:
//   - task: The task to evaluate.
//   - now: The current timestamp.
//
// Returns:
//   - []TaskModifier: A list of updates to apply to the task.
func getTaskUpdates(task models.Task, now time.Time) []TaskModifier {
	updates := make([]TaskModifier, 0, 2)
	if task.Content.IsValid() && !task.IsProject {
		updates = append(updates, ProjectModifier)
	}
	if task.Done {
		updates = append(updates, CompletionModifier(now))
	}
	return updates
}

// DetermineTaskAction decides the appropriate action for a task based on its state and retention policy.
//
// Parameters:
//   - task: The task to evaluate.
//   - now: The current timestamp.
//   - config: Retention rules for completed tasks.
//
// Returns:
//   - TaskAction: The determined action and any required updates.
func DetermineTaskAction(task models.Task, now time.Time, config RetentionConfig) TaskAction {
	if IsCompleted(task) {
		if !ShouldRetainTask(task, now, config) {
			return TaskAction{
				Task:   task,
				Action: ActionArchive,
			}
		}
		return TaskAction{
			Task:   task,
			Action: ActionRetain,
		}
	}

	return TaskAction{
		Task:    task,
		Action:  ActionUpdate,
		Updates: getTaskUpdates(task, now),
	}
}

// ApplyTaskAction applies the specified updates to a task and returns the updated task.
//
// Parameters:
//   - action: The action containing the updates to apply.
//   - now: The current timestamp.
//
// Returns:
//   - models.Task: The updated task.
func ApplyTaskAction(action TaskAction, now time.Time) models.Task {
	if len(action.Updates) == 0 {
		return action.Task
	}
	return ComposeModifiers(action.Updates...)(action.Task, now)
}

// ProcessTasks processes a batch of tasks, determining and optionally applying actions for each.
//
// Parameters:
//   - tasks: The list of tasks to process.
//   - now: The current timestamp.
//   - config: Retention rules for completed tasks.
//
// Returns:
//   - ptr.Option[[]TaskAction]: A list of actions for the tasks, or None if the list is empty.
func ProcessTasks(tasks []models.Task, now time.Time, config RetentionConfig) ptr.Option[[]TaskAction] {
	if len(tasks) == 0 {
		return ptr.None[[]TaskAction]()
	}

	processor := func(task models.Task) TaskAction {
		action := DetermineTaskAction(task, now, config)
		if action.Action == ActionUpdate {
			action.Task = ApplyTaskAction(action, now)
		}
		return action
	}

	return ptr.Some(Map(tasks, processor))
}

// handleTaskAction executes the specified action on a task (e.g., retain, archive, update, complete).
// Currently, archive, update, and complete actions are placeholders.
//
// Parameters:
//   - action: The action to execute.
//   - basePath: Base path for storage operations (if applicable).
//   - now: The current timestamp.
//
// Returns:
//   - error: An error if the action cannot be executed or is unsupported.
func handleTaskAction(action TaskAction, basePath string, now time.Time) error {
	_ = basePath
	_ = now

	switch action.Action {
	case ActionRetain:
		// For retain, we don't need to do anything
		return nil

	case ActionArchive:
		// TODO: Implement archive logic
		return fmt.Errorf("archive action not implemented")

	case ActionUpdate:
		// TODO: Implement update logic
		return fmt.Errorf("update action not implemented")

	case ActionComplete:
		// TODO: Implement complete logic
		return fmt.Errorf("complete action not implemented")

	default:
		return fmt.Errorf("unknown action type")
	}
}
