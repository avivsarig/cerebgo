package tasks

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func TestIsCompleted(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	completedTime := baseTime.Add(1 * time.Hour)
	baseDateStr := baseTime.Format("2006-01-02")

	baseTask := models.Task{
		Title:     "Test Task",
		Content:   ptr.None[string](),
		CreatedAt: baseTime,
		DoDate:    baseDateStr,
	}

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "neither_flag_set",
			task: models.Task{
				Title:       baseTask.Title,
				Content:     baseTask.Content,
				Done:        false,
				CompletedAt: ptr.None[time.Time](),
				CreatedAt:   baseTask.CreatedAt,
				DoDate:      baseTask.DoDate,
			},
			want: false,
		},
		{
			name: "only_done_set",
			task: models.Task{
				Title:       baseTask.Title,
				Content:     baseTask.Content,
				Done:        true,
				CompletedAt: ptr.None[time.Time](),
				CreatedAt:   baseTask.CreatedAt,
				DoDate:      baseTask.DoDate,
			},
			want: false,
		},
		{
			name: "only_completed_at_set",
			task: models.Task{
				Title:       baseTask.Title,
				Content:     baseTask.Content,
				Done:        false,
				CompletedAt: ptr.Some(completedTime),
				CreatedAt:   baseTask.CreatedAt,
				DoDate:      baseTask.DoDate,
			},
			want: false,
		},
		{
			name: "both_flags_set",
			task: models.Task{
				Title:       baseTask.Title,
				Content:     baseTask.Content,
				Done:        true,
				CompletedAt: ptr.Some(completedTime),
				CreatedAt:   baseTask.CreatedAt,
				DoDate:      baseTask.DoDate,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCompleted(tt.task)
			if got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}
