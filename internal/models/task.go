package models

import (
	"time"

	"github.com/avivSarig/cerebgo/pkg/ptr"
)

type Task struct {
	Title          string
	Content        ptr.Option[string]
	IsProject      bool
	IsHighPriority bool
	Done           bool
	CompletedAt    ptr.Option[time.Time]
	DueDate        ptr.Option[string] // YYYY-MM-DD format
	DoDate         string             // YYYY-MM-DD format, required
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
