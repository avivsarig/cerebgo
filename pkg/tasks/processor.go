package tasks

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func ProcessAllTasks(now time.Time, configuration *viper.Viper) error {
	retentionConfig := RetentionConfig{
		EmptyTaskRetention: time.Duration(configuration.GetInt("settings.retention.empty_task")) * 24 * time.Hour,
		ProjectRetention:   time.Duration(configuration.GetInt("settings.retention.project_before_archive")) * 24 * time.Hour,
	}

	activeTasksPath := filepath.Join(
		configuration.GetString("base_path"),
		configuration.GetString("paths.base.tasks"),
	)
	completedTasksPath := filepath.Join(
		configuration.GetString("base_path"),
		configuration.GetString("paths.subdirs.tasks.completed"),
	)
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

		resultTask, err := ApplyModifiers(task, now, modifiers...)
		if err != nil {
			return fmt.Errorf("failed to apply task modifiers: %w", err)
		}

		err = RewriteTask(resultTask, completedTasksPath)
		if err != nil {
			return fmt.Errorf("failed to rewrite task: %w", err)
		}
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

		resultTask, err := ApplyModifiers(task, now, modifiers...)
		if err != nil {
			return fmt.Errorf("failed to apply task modifiers: %w", err)
		}

		err = RewriteTask(resultTask, activeTasksPath)
		if err != nil {
			return fmt.Errorf("failed to rewrite task: %w", err)
		}
	}

	return nil
}
