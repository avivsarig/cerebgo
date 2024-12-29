package mdparser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// MarkdownDocument represents a parsed markdown content with its metadata.
type MarkdownDocument struct {
	Title       string
	Frontmatter map[string]interface{}
	Content     string
}

// ParseMarkdownDoc parses raw markdown text into a structured MarkdownDocument.
func ParseMarkdownDoc(filePath string) (MarkdownDocument, error) {
	baseFile := filepath.Base(filePath)
	title := strings.TrimSuffix(baseFile, filepath.Ext(baseFile))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return MarkdownDocument{}, fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)

	// Empty file is valid
	if strings.TrimSpace(content) == "" {
		return MarkdownDocument{Title: title}, nil
	}

	// First check whitespace
	if strings.TrimLeft(content, " \t\n") != content {
		return MarkdownDocument{}, fmt.Errorf("invalid markdown: whitespace before frontmatter")
	}

	// Check for frontmatter markers
	trimmedContent := strings.TrimSpace(content)
	if strings.HasPrefix(trimmedContent, "--") {
		// If it starts with any number of dashes but not exactly 3, it's invalid
		if !strings.HasPrefix(trimmedContent, "---") || strings.HasPrefix(trimmedContent, "----") {
			return MarkdownDocument{}, fmt.Errorf("invalid markdown: incorrect frontmatter markers")
		}
	} else if strings.Contains(content, "---") {
		return MarkdownDocument{}, fmt.Errorf("invalid markdown: content before frontmatter")
	} else {
		return MarkdownDocument{Title: title, Content: trimmedContent}, nil
	}

	// Look for closing frontmatter marker
	lines := strings.Split(content, "\n")
	var closingIndex int
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			closingIndex = i
			break
		}
	}

	if closingIndex == 0 {
		return MarkdownDocument{}, fmt.Errorf("invalid markdown: unclosed frontmatter")
	}

	// Check for multiple frontmatter blocks
	firstClose := strings.Index(content[3:], "\n---\n")
	if firstClose != -1 {
		secondBlock := strings.Index(content[firstClose+7:], "\n---\n")
		if secondBlock != -1 {
			return MarkdownDocument{}, fmt.Errorf("invalid markdown: multiple frontmatter blocks")
		}
	}

	// Parse frontmatter
	frontmatter := strings.Join(lines[1:closingIndex], "\n")
	fm := make(map[string]interface{})

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(frontmatter)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&fm); err != nil {
		return MarkdownDocument{}, fmt.Errorf("invalid frontmatter YAML: %w", err)
	}

	// Convert time values to strings
	for k, v := range fm {
		if t, ok := v.(time.Time); ok {
			fm[k] = t.Format("2006-01-02")
		}
	}

	remainingContent := ""
	if closingIndex+1 < len(lines) {
		remainingContent = strings.TrimSpace(strings.Join(lines[closingIndex+1:], "\n"))
	}

	return MarkdownDocument{
		Title:       title,
		Frontmatter: fm,
		Content:     remainingContent,
	}, nil
}
