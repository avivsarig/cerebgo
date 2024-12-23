package mdparser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// WriteMarkdownDoc writes a markdown document with frontmatter to a file.
// The frontmatter is converted to YAML format and enclosed in --- markers.
//
// Parameters:
//   - fm: Frontmatter metadata as key-value pairs
//   - content: Main markdown content
//   - path: File path to write the document to
//
// Returns:
//   - error if marshaling frontmatter or writing file fails
func WriteMarkdownDoc(fm Frontmatter, content string, path string) error {
	// create frontmatter string
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	// construct markdown document
	mdDoc := fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), content)

	// write file to path
	err = os.WriteFile(path, []byte(mdDoc), 0644)
	if err != nil {
		return fmt.Errorf("failed to write markdown document to %s: %w", path, err)
	}
	return nil
}
