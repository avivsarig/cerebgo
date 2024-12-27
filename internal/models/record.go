package models

import (
	"time"

	"github.com/avivSarig/cerebgo/pkg/ptr"
)

type Record struct {
	Title      string
	Content    ptr.Option[string]
	Tags       []string
	URL        ptr.Option[string]
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt ptr.Option[time.Time]
}
