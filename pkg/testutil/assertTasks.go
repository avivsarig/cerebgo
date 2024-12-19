package testutil

import (
	"testing"

	"github.com/avivSarig/cerebgo/internal/models"
)

func AssertTaskEqual(t *testing.T, got, want models.Task) {
	t.Helper()

	results := []ValidationResult{
		// Simple fields
		ValidateEqual("Title", got.Title, want.Title),
		ValidateEqual("IsProject", got.IsProject, want.IsProject),
		ValidateEqual("IsHighPriority", got.IsHighPriority, want.IsHighPriority),
		ValidateEqual("Done", got.Done, want.Done),

		// Optional fields
		ValidateOptional("Content", got.Content, want.Content, StringComparer),
		ValidateOptional("CompletedAt", got.CompletedAt, want.CompletedAt, TimeComparer),
		ValidateOptional("DueDate", got.DueDate, want.DueDate, StringComparer),

		// Required fields
		ValidateEqual("DoDate", got.DoDate, want.DoDate),
		ValidateEqual("CreatedAt", got.CreatedAt, want.CreatedAt),
		ValidateEqual("UpdatedAt", got.UpdatedAt, want.UpdatedAt),
	}

	ReportResults(t, results)
}
