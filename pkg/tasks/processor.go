package tasks

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// Access values.
func GetConfigValue() string {
	return viper.GetString("key.name")
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
