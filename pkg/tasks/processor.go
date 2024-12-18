package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/spf13/viper"
)

type Processor struct {
	config *viper.Viper
}

func NewProcessor(config *viper.Viper) (*Processor, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Processor{
		config: config,
	}, nil
}

func (p *Processor) ClearCompletedTasks(now time.Time) error {
	completedPath := p.config.GetString("paths.subdirs.tasks.completed")
	emptyTaskRetention := p.config.GetInt("settings.retention.empty_task")
	projectRetention := p.config.GetInt("settings.retention.project_before_archive")

	entries, err := os.ReadDir(completedPath)
	if err != nil {
		return fmt.Errorf("failed to read completed tasks directory: %w", err)
	}

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(completedPath, entry.Name())
		doc, err := mdparser.ParseMarkdownDoc(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse task file %s: %w", entry.Name(), err)
		}

		task, err := mdparser.DocumetToTask(doc)
		if err != nil {
			return fmt.Errorf("failed to convert document to task %s: %w", entry.Name(), err)
		}

		// TODO: !task.Done
		if !task.CompletedAt.IsValid() {
			// TODO: return to active tasks!
			continue
		}

		completedAge := now.Sub(task.CompletedAt.Value())

		shouldDelete := false
		if task.IsProject {
			shouldDelete = completedAge > time.Duration(projectRetention)*24*time.Hour
		} else {
			shouldDelete = completedAge > time.Duration(emptyTaskRetention)*24*time.Hour
		}

		if shouldDelete {
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to delete task file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

func ProcessTasks() error {
	// Clear ${paths.subdirs.tasks.completed} from tasks:
	// if `is_project=true` and `completed_at`>${settings.retention.project_before_archive} ago
	// if `is_project=false` and `completed_at`>${settings.retention.empty_task}

	// Parse and create new tasks from:
	// - ${paths.base.journal}
	// - ${paths.base.inbox}

	// Process tasks from ${paths.base.tasks}:
	// 1. if content!='' --> `is_project=true`
	// 2. if `done=true` --> `completed_at`=now() --> move to ${paths.subdirs.tasks.completed}
	// 3. if `do_date`<date(now()) --> `do_date`=date(now)

	return nil
}
