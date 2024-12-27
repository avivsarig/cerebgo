package mdparser

import (
	"fmt"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func DocumentToTask(doc MarkdownDocument) (models.Task, error) {
	// Helper to reduce error handling boilerplate
	getFrontmatter := func(field string) (string, error) {
		value, err := getFrontmatterString(doc.Frontmatter, field)
		if err != nil {
			return "", fmt.Errorf("getting %s: %w", field, err)
		}
		return value, nil
	}

	// Get and parse required fields
	createdAtStr, err := getFrontmatter("created_at")
	if err != nil {
		return models.Task{}, err
	}
	createdAt, err := parseDate(createdAtStr)
	if err != nil {
		return models.Task{}, fmt.Errorf("parsing created_at: %w", err)
	}

	updatedAtStr, err := getFrontmatter("updated_at")
	if err != nil {
		return models.Task{}, err
	}
	updatedAt, err := parseDate(updatedAtStr)
	if err != nil {
		return models.Task{}, fmt.Errorf("parsing updated_at: %w", err)
	}

	// Handle optional due date
	var dueDateOpt ptr.Option[string]
	if dueDateStr, err := getFrontmatter("due_date"); err == nil {
		dueDate, err := parseDate(dueDateStr)
		if err != nil {
			return models.Task{}, fmt.Errorf("parsing due_date: %w", err)
		}
		dueDateOpt = ptr.Some(dueDate.Format(time.RFC3339))
	} else {
		dueDateOpt = ptr.None[string]()
	}

	// Handle priority with proper error checking
	priority, err := getFrontmatter("priority")
	if err != nil {
		// Assuming no priority means not high priority
		priority = "normal"
	}
	isHighPriority := priority == "high"

	// Validate required fields
	if doc.Title == "" {
		return models.Task{}, fmt.Errorf("task title is required")
	}

	is_project := doc.Content != ""

	return models.Task{
		Title:          doc.Title,
		Content:        ptr.Some(doc.Content),
		IsProject:      is_project,
		IsHighPriority: isHighPriority,
		Done:           false,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DueDate:        dueDateOpt,
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
