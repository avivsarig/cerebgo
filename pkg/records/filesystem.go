package records

import (
	"path/filepath"
	"time"

	"github.com/avivSarig/cerebgo/internal/models"
	"github.com/avivSarig/cerebgo/pkg/mdparser"
)

// WriteRecordToFile writes a Record to a markdown file
//
// Parameters:
//   - record: Record to write
//
// Returns:
//   - error: writing file errors with context
//
// FUTURE: consider add overwrite flag (at the moment, it always overwrites).
func WriteRecordToFile(record models.Record, path string) error {
	fm := mdparser.Frontmatter{
		"tags":       record.Tags,
		"created_at": record.CreatedAt,
		"updated_at": record.UpdatedAt,
	}

	if record.URL.IsValid() {
		fm["url"] = record.URL.Value()
	}

	if record.ArchivedAt.IsValid() {
		fm["archived_at"] = record.ArchivedAt.Value()
	} else {
		fm["archived_at"] = time.Now().Format(time.RFC3339)
	}

	content := ""
	if record.Content.IsValid() {
		content = record.Content.Value()
	}

	filename := filepath.Join(path, record.Title+".md")
	return mdparser.WriteMarkdownDoc(fm, content, filename)
}
