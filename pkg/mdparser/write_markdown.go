package mdparser

import (
	"fmt"
	"os"
	"strings"

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
	// Validate no function values in frontmatter
	for _, v := range fm {
		if vType := fmt.Sprintf("%T", v); strings.Contains(vType, "func(") {
			return fmt.Errorf("frontmatter contains unsupported function value")
		}
	}

	var fmBytes []byte
	var err error

	if len(fm) > 0 {
		fmBytes, err = yaml.Marshal(fm)
		if err != nil {
			return fmt.Errorf("failed to marshal frontmatter: %w", err)
		}
	}

	mdDoc := fmt.Sprintf("---%s---\n\n%s",
		func() string {
			if len(fmBytes) > 0 {
				return "\n" + string(fmBytes)
			}
			return "\n"
		}(),
		content)

	return os.WriteFile(path, []byte(mdDoc), 0644)
}
