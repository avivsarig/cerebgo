package tasks

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Processor struct {
	config *viper.Viper
}

func NewProcessor(config *viper.Viper) (*Processor, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &Processor{config: config}, nil
}

func (p *Processor) ProcessAllTasks(now time.Time) error {
	config := RetentionConfig{
		EmptyTaskRetention: time.Duration(p.config.GetInt("settings.retention.empty_task")) * 24 * time.Hour,
		ProjectRetention:   time.Duration(p.config.GetInt("settings.retention.project_before_archive")) * 24 * time.Hour,
	}

	// TODO: CLEAR COMPLETED TASKS FIRST
	tasks, err := readTasksFromDirectory(p.config.GetString("paths.base.tasks"))
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	result := ProcessTasks(tasks, now, config)
	if !result.IsValid() {
		return fmt.Errorf("task processing failed")
	}

	// Process the results...
	return nil
}
