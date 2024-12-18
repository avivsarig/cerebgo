package tasks

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestCompletionModifier(t *testing.T) {
	// We'll fix our reference time to make tests deterministic
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	baseDateStr := baseTime.Format("2006-01-02")
	futureDateStr := baseTime.Add(24 * time.Hour).Format("2006-01-02")

	tests := []struct {
		name           string
		input          models.Task
		completionTime time.Time
		currentTime    time.Time
		want           models.Task
	}{
		{
			name: "complete_uncompleted_task",
			input: models.Task{
				Title:          "Test Task",
				Content:        ptr.Some("content"),
				IsProject:      false,
				IsHighPriority: true,
				CompletedAt:    ptr.None[time.Time](),
				CreatedAt:      baseTime.Add(-2 * time.Hour),
				DueDate:        ptr.Some(futureDateStr),
				DoDate:         baseDateStr,
			},
			completionTime: baseTime.Add(-1 * time.Hour),
			currentTime:    baseTime,
			want: models.Task{
				Title:          "Test Task",
				Content:        ptr.Some("content"),
				IsProject:      false,
				IsHighPriority: true,
				CompletedAt:    ptr.Some(baseTime.Add(-1 * time.Hour)),
				CreatedAt:      baseTime.Add(-2 * time.Hour),
				UpdatedAt:      baseTime,
				DueDate:        ptr.Some(futureDateStr),
				DoDate:         baseDateStr,
			},
		},
		{
			name: "complete_already_completed_task",
			input: models.Task{
				Title:          "Already Complete Task",
				Content:        ptr.Some("content"),
				CompletedAt:    ptr.Some(baseTime.Add(-2 * time.Hour)),
				CreatedAt:      baseTime.Add(-3 * time.Hour),
				IsHighPriority: false,
			},
			completionTime: baseTime.Add(-1 * time.Hour),
			currentTime:    baseTime,
			want: models.Task{
				Title:          "Already Complete Task",
				Content:        ptr.Some("content"),
				CompletedAt:    ptr.Some(baseTime.Add(-1 * time.Hour)), // Should override previous completion time
				CreatedAt:      baseTime.Add(-3 * time.Hour),
				UpdatedAt:      baseTime,
				IsHighPriority: false,
			},
		},
		{
			name: "complete_task_with_null_content",
			input: models.Task{
				Title:       "Minimal Task",
				Content:     ptr.None[string](),
				CompletedAt: ptr.None[time.Time](),
				CreatedAt:   baseTime.Add(-1 * time.Hour),
			},
			completionTime: baseTime,
			currentTime:    baseTime,
			want: models.Task{
				Title:       "Minimal Task",
				Content:     ptr.None[string](),
				CompletedAt: ptr.Some(baseTime),
				CreatedAt:   baseTime.Add(-1 * time.Hour),
				UpdatedAt:   baseTime,
			},
		},
		{
			name: "complete_project_task",
			input: models.Task{
				Title:       "Project Task",
				Content:     ptr.Some("project content"),
				IsProject:   true,
				CompletedAt: ptr.None[time.Time](),
				CreatedAt:   baseTime.Add(-1 * time.Hour),
			},
			completionTime: baseTime,
			currentTime:    baseTime,
			want: models.Task{
				Title:       "Project Task",
				Content:     ptr.Some("project content"),
				IsProject:   true,
				CompletedAt: ptr.Some(baseTime),
				CreatedAt:   baseTime.Add(-1 * time.Hour),
				UpdatedAt:   baseTime,
			},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompletionModifier(tt.completionTime)(tt.input, tt.currentTime)
			testutil.AssertTaskEqual(t, got, tt.want)
		})
	}
}

func TestProjectModifier(t *testing.T) {
	// Fix reference time for deterministic tests
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	baseDateStr := baseTime.Format("2006-01-02")
	futureDateStr := baseTime.Add(24 * time.Hour).Format("2006-01-02")

	tests := []struct {
		name        string
		input       models.Task
		currentTime time.Time
		want        models.Task
	}{
		{
			name: "convert_task_to_project",
			input: models.Task{
				Title:          "Regular Task",
				Content:        ptr.Some("content"),
				IsProject:      false,
				IsHighPriority: true,
				CompletedAt:    ptr.None[time.Time](),
				CreatedAt:      baseTime.Add(-1 * time.Hour),
				DueDate:        ptr.Some(futureDateStr),
				DoDate:         baseDateStr,
			},
			currentTime: baseTime,
			want: models.Task{
				Title:          "Regular Task",
				Content:        ptr.Some("content"),
				IsProject:      true, // Should be converted to project
				IsHighPriority: true,
				CompletedAt:    ptr.None[time.Time](),
				CreatedAt:      baseTime.Add(-1 * time.Hour),
				UpdatedAt:      baseTime,
				DueDate:        ptr.Some(futureDateStr),
				DoDate:         baseDateStr,
			},
		},
		{
			name: "already_project_task",
			input: models.Task{
				Title:          "Existing Project",
				Content:        ptr.Some("project content"),
				IsProject:      true,
				IsHighPriority: false,
				CompletedAt:    ptr.Some(baseTime.Add(-2 * time.Hour)),
				CreatedAt:      baseTime.Add(-3 * time.Hour),
			},
			currentTime: baseTime,
			want: models.Task{
				Title:          "Existing Project",
				Content:        ptr.Some("project content"),
				IsProject:      true,
				IsHighPriority: false,
				CompletedAt:    ptr.Some(baseTime.Add(-2 * time.Hour)),
				CreatedAt:      baseTime.Add(-3 * time.Hour),
				UpdatedAt:      baseTime,
			},
		},
		{
			name: "minimal_task_to_project",
			input: models.Task{
				Title:       "Minimal Task",
				Content:     ptr.None[string](),
				CreatedAt:   baseTime.Add(-1 * time.Hour),
				CompletedAt: ptr.None[time.Time](),
			},
			currentTime: baseTime,
			want: models.Task{
				Title:       "Minimal Task",
				Content:     ptr.None[string](),
				IsProject:   true,
				CreatedAt:   baseTime.Add(-1 * time.Hour),
				UpdatedAt:   baseTime,
				CompletedAt: ptr.None[time.Time](),
			},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProjectModifier(tt.input, tt.currentTime)
			testutil.AssertTaskEqual(t, got, tt.want)
		})
	}
}

func TestComposeModifiers(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	baseDateStr := baseTime.Format("2006-01-02")

	tests := []struct {
		name    string
		input   models.Task
		current time.Time
		// List the modifiers we want to compose
		modifiers []TaskModifier
		want      models.Task
	}{
		{
			name: "no_modifiers",
			input: models.Task{
				Title:     "Basic Task",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-1 * time.Hour),
			},
			current:   baseTime,
			modifiers: []TaskModifier{}, // Empty modifier list
			want: models.Task{
				Title:     "Basic Task",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-1 * time.Hour),
			},
		},
		{
			name: "single_modifier",
			input: models.Task{
				Title:     "Task to Project",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-1 * time.Hour),
			},
			current: baseTime,
			modifiers: []TaskModifier{
				ProjectModifier,
			},
			want: models.Task{
				Title:     "Task to Project",
				Content:   ptr.Some("content"),
				IsProject: true,
				CreatedAt: baseTime.Add(-1 * time.Hour),
				UpdatedAt: baseTime,
			},
		},
		{
			name: "multiple_modifiers_order",
			input: models.Task{
				Title:     "Complex Task",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-2 * time.Hour),
				DoDate:    baseDateStr,
			},
			current: baseTime,
			modifiers: []TaskModifier{
				ProjectModifier,
				CompletionModifier(baseTime.Add(-1 * time.Hour)),
			},
			want: models.Task{
				Title:       "Complex Task",
				Content:     ptr.Some("content"),
				IsProject:   true,
				CreatedAt:   baseTime.Add(-2 * time.Hour),
				UpdatedAt:   baseTime,
				CompletedAt: ptr.Some(baseTime.Add(-1 * time.Hour)),
				DoDate:      baseDateStr,
			},
		},
		{
			name: "multiple_modifiers_reverse_order",
			input: models.Task{
				Title:     "Complex Task Reverse",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-2 * time.Hour),
				DoDate:    baseDateStr,
			},
			current: baseTime,
			modifiers: []TaskModifier{
				CompletionModifier(baseTime.Add(-1 * time.Hour)),
				ProjectModifier,
			},
			want: models.Task{
				Title:       "Complex Task Reverse",
				Content:     ptr.Some("content"),
				IsProject:   true,
				CreatedAt:   baseTime.Add(-2 * time.Hour),
				UpdatedAt:   baseTime,
				CompletedAt: ptr.Some(baseTime.Add(-1 * time.Hour)),
				DoDate:      baseDateStr,
			},
		},
		{
			name: "idempotent_modifiers",
			input: models.Task{
				Title:     "Idempotent Test",
				Content:   ptr.Some("content"),
				CreatedAt: baseTime.Add(-2 * time.Hour),
			},
			current: baseTime,
			modifiers: []TaskModifier{
				ProjectModifier,
				ProjectModifier, // Apply same modifier twice
			},
			want: models.Task{
				Title:     "Idempotent Test",
				Content:   ptr.Some("content"),
				IsProject: true,
				CreatedAt: baseTime.Add(-2 * time.Hour),
				UpdatedAt: baseTime,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composedModifier := ComposeModifiers(tt.modifiers...)
			got := composedModifier(tt.input, tt.current)
			testutil.AssertTaskEqual(t, got, tt.want)
		})
	}
}
