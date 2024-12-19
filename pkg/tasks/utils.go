package tasks

import "github.com/avivSarig/cerebgo/internal/models"

func IsCompleted(t models.Task) bool {
	return t.Done && t.CompletedAt.IsValid()
}
