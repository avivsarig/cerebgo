package files

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func MoveFile(srcPath, destPath, fileName string) error {
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if !slices.ContainsFunc(entries, func(e os.DirEntry) bool {
		return e.Name() == fileName
	}) {
		return fmt.Errorf("file %s not found in %s", fileName, srcPath)
	}

	oldPath := filepath.Join(srcPath, fileName)
	newPath := filepath.Join(destPath, fileName)
	return os.Rename(oldPath, newPath)
}
