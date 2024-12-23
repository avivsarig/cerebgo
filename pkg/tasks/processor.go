package tasks

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func ProcessAllTasks(now time.Time, configuration *viper.Viper) error {
	retentionConfig := RetentionConfig{
		EmptyTaskRetention: time.Duration(configuration.GetInt("settings.retention.empty_task")) * 24 * time.Hour,
		ProjectRetention:   time.Duration(configuration.GetInt("settings.retention.project_before_archive")) * 24 * time.Hour,
	}

	completedTasksPath := configuration.GetString("paths.subdir.tasks.completed")
	activeTasksPath := configuration.GetString("paths.subdir.tasks.completed")

	// TODO: add logging

	// process completed tasks:
	completedTasks, err := readTasksFromDirectory(completedTasksPath)
	if err != nil {
		return fmt.Errorf("failed to read completed tasks: %w", err)
	}
	for _, task := range completedTasks {
		modifiers, err := PlanCompletedTaskActions(task, time.Now(), retentionConfig)
		if err != nil {
			return fmt.Errorf("failed to process completed tasks: %w", err)
		}

		// TODO: remove placeholder
		_ = modifiers
		// for _, modifier := range modifiers {
		// 	task, err = modifier(task)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to apply task modifier: %w", err)
		// 	}
		// }
	}

	// process active tasks:
	activeTasks, err := readTasksFromDirectory(activeTasksPath)
	if err != nil {
		return fmt.Errorf("failed to read active tasks: %w", err)
	}

	for _, task := range activeTasks {
		modifiers, err := PlanActiveTaskActions(task, time.Now())
		if err != nil {
			return fmt.Errorf("failed to process active tasks: %w", err)
		}

		// TODO: remove placeholder
		_ = modifiers
		// for _, modifier := range modifiers {
		// 	task, err = modifier(task)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to apply task modifier: %w", err)
		// 	}
		// }
	}

	return nil
}
