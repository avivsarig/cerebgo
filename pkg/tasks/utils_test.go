package tasks_test

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/tasks"
)

func TestIsCompleted(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "done with completed_at",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.Some(now),
			},
			want: true,
		},
		{
			name: "done without completed_at",
			task: models.Task{
				Done:        true,
				CompletedAt: ptr.None[time.Time](),
			},
			want: false,
		},
		{
			name: "not done with completed_at",
			task: models.Task{
				Done:        false,
				CompletedAt: ptr.Some(now),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tasks.IsCompleted(tt.task); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDoDate(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		task models.Task
		now  time.Time
		want bool
	}{
		{
			name: "future date",
			task: models.Task{DoDate: "2024-12-31"},
			now:  fixedTime,
			want: true,
		},
		{
			name: "today's date",
			task: models.Task{DoDate: "2024-01-01"},
			now:  fixedTime,
			want: true,
		},
		{
			name: "past date",
			task: models.Task{DoDate: "2023-12-31"},
			now:  fixedTime,
			want: false,
		},
		{
			name: "invalid format",
			task: models.Task{DoDate: "01-01-2024"},
			now:  fixedTime,
			want: false,
		},
		{
			name: "empty date",
			task: models.Task{DoDate: ""},
			now:  fixedTime,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tasks.IsValidDoDate(tt.task, tt.now); got != tt.want {
				t.Errorf("IsValidDoDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDueDate(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		task models.Task
		now  time.Time
		want bool
	}{
		{
			name: "future date",
			task: models.Task{DueDate: ptr.Some("2024-12-31")},
			now:  fixedTime,
			want: true,
		},
		{
			name: "today's date",
			task: models.Task{DueDate: ptr.Some("2024-01-01")},
			now:  fixedTime,
			want: true,
		},
		{
			name: "past date",
			task: models.Task{DueDate: ptr.Some("2023-12-31")},
			now:  fixedTime,
			want: false,
		},
		{
			name: "no due date",
			task: models.Task{DueDate: ptr.None[string]()},
			now:  fixedTime,
			want: false,
		},
		{
			name: "invalid format",
			task: models.Task{DueDate: ptr.Some("01-01-2024")},
			now:  fixedTime,
			want: false,
		},
		{
			name: "empty date",
			task: models.Task{DueDate: ptr.Some("")},
			now:  fixedTime,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tasks.IsValidDueDate(tt.task, tt.now); got != tt.want {
				t.Errorf("IsValidDueDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
