package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// FilePath represents a file path.
// It contains the directory and the name of the file.
//
// The full path can be obtained by calling FullPath().
type FilePath struct {
	Dir  string
	Name string
}

// FullPath returns the full path of the file.
//
// Returns:
//   - string: the full path of the file.
func (f FilePath) FullPath() string {
	return filepath.Join(f.Dir, f.Name)
}

// FileExists checks if a file exists at the specified path.
//
// Parameters:
//   - path: path to the file to check.
//
// Returns:
//   - bool: true if the file exists, false otherwise.
func FileExists(path FilePath) (bool, error) {
	entries, err := os.ReadDir(path.Dir)
	if err != nil {
		return false, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == path.Name {
			return true, nil
		}
	}
	return false, nil
}

// MoveFile moves a file from source to destination path. Uses os.Rename which is atomic on POSIX systems
//
// Parameters:
//   - src: source file path
//   - dest: destination file path
//
// Returns error if:
//   - source file doesn't exist
//   - directory read fails
//   - rename operation fails
func MoveFile(src, dest FilePath) error {
	exists, err := FileExists(src)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file %s not found in %s", src.Name, src.Dir)
	}

	return os.Rename(src.FullPath(), dest.FullPath())
}

// DeleteFile removes a file at the specified path.
//
// Parameters:
//   - src: path to the file to delete.
//
// Returns error if file doesn't exist or deletion fails.
func DeleteFile(src FilePath) error {
	exists, err := FileExists(src)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file %s not found in %s", src.Name, src.Dir)
	}

	return os.Remove(src.FullPath())
}
