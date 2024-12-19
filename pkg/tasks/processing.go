package tasks

import (
	"fmt"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

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

func ApplyTaskAction(action TaskAction, now time.Time) models.Task {
	if len(action.Updates) == 0 {
		return action.Task
	}
	return ComposeModifiers(action.Updates...)(action.Task, now)
}

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
