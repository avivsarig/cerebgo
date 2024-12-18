package mdparser

import (
	"fmt"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func DocumetToTask(doc MarkdownDocument) (models.Task, error) {
	createdAtStr, err := getFrontmatterString(doc.Frontmatter, "created_at")
	if err != nil {
		return models.Task{}, fmt.Errorf("getting created_at: %w", err)
	}

	createdAt, err := parseDate(createdAtStr)
	if err != nil {
		return models.Task{}, fmt.Errorf("parsing created_at: %w", err)
	}

	updatedAtStr, err := getFrontmatterString(doc.Frontmatter, "updated_at")
	if err != nil {
		return models.Task{}, fmt.Errorf("getting updated_at: %w", err)
	}

	updatedAt, err := parseDate(updatedAtStr)
	if err != nil {
		return models.Task{}, fmt.Errorf("parsing updated_at: %w", err)
	}

	// TODO: FIX
	// // Get optional fields
	// dueDate, err := getFrontmatterString(doc.Frontmatter, "due_date")
	// var dueDateStr ptr.Option[string]
	// if dueDate.IsValid() {
	// 	// Convert time.Time back to string
	// 	dueDateStr = ptr.Some(dueDate.Value().Format(time.RFC3339))
	// } else {
	// 	dueDateStr = ptr.None[string]()
	// }

	// Get priority field and convert to boolean
	// priority, err := getFrontmatterString(doc.Frontmatter, "priority")
	// isHighPriority := priority == "high"

	return models.Task{
		Title:          doc.Title,
		Content:        ptr.Some(doc.Content),
		IsProject:      false,
		IsHighPriority: false,
		Done:           false,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DueDate:        ptr.None[string](),
		CompletedAt:    ptr.None[time.Time](),
	}, nil
}
func getFrontmatterString(fm map[string]interface{}, key string) (string, error) {
	value, exists := fm[key]
	if !exists {
		return "", fmt.Errorf("field %s not found", key)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not a string", key)
	}

	return strValue, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}
