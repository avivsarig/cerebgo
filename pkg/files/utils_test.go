package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/avivSarig/cerebgo/pkg/files"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// TestFilePath_FullPath tests the FullPath method of FilePath struct.
// It verifies that paths are correctly joined according to the OS-specific path separator.
func TestFilePath_FullPath(t *testing.T) {
	tests := []struct {
		name string
		path files.FilePath
		want string
	}{
		{
			name: "simple path",
			path: files.FilePath{Dir: "testdir", Name: "test.txt"},
			want: filepath.Join("testdir", "test.txt"),
		},
		{
			name: "path with multiple segments",
			path: files.FilePath{Dir: "path/to/dir", Name: "file.txt"},
			want: filepath.Join("path", "to", "dir", "file.txt"),
		},
		{
			name: "root directory",
			path: files.FilePath{Dir: "/", Name: "root.txt"},
			want: filepath.Join("/", "root.txt"),
		},
		{
			name: "empty directory",
			path: files.FilePath{Dir: "", Name: "file.txt"},
			want: "file.txt",
		},
		{
			name: "both empty",
			path: files.FilePath{Dir: "", Name: ""},
			want: "",
		},
		{
			name: "directory with dots",
			path: files.FilePath{Dir: "../test/./dir", Name: "file.txt"},
			want: filepath.Join("..", "test", "dir", "file.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.FullPath()
			if got != tt.want {
				t.Errorf("FullPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFileExists tests the FileExists function.
// It creates temporary test files and verifies existence checks work correctly.
func TestFileExists(t *testing.T) {
	// Create a temporary test directory
	testDir := testutil.CreateTestDirectory(t)

	// Create a test file
	testFile := "test.txt"
	err := testutil.CreateTestFile(t, testDir, testFile, "test content")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    files.FilePath
		want    bool
		wantErr bool
	}{
		{
			name: "existing file",
			path: files.FilePath{Dir: testDir, Name: testFile},
			want: true,
		},
		{
			name: "non-existent file",
			path: files.FilePath{Dir: testDir, Name: "nonexistent.txt"},
			want: false,
		},
		{
			name:    "invalid directory",
			path:    files.FilePath{Dir: "/nonexistent/dir", Name: "file.txt"},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := files.FileExists(tt.path)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, don't check the return value
			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMoveFile tests the MoveFile function.
// It verifies file movement operations work correctly in various scenarios.
func TestMoveFile(t *testing.T) {
	// Create a temporary test directory
	testDir := testutil.CreateTestDirectory(t)
	destDir := filepath.Join(testDir, "dest")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Create a test file
	testFile := "source.txt"
	testContent := "test content"
	err := testutil.CreateTestFile(t, testDir, testFile, testContent)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		src     files.FilePath
		dest    files.FilePath
		wantErr bool
	}{
		{
			name: "valid move",
			src:  files.FilePath{Dir: testDir, Name: testFile},
			dest: files.FilePath{Dir: destDir, Name: "dest.txt"},
		},
		{
			name:    "source doesn't exist",
			src:     files.FilePath{Dir: testDir, Name: "nonexistent.txt"},
			dest:    files.FilePath{Dir: destDir, Name: "dest2.txt"},
			wantErr: true,
		},
		{
			name:    "invalid source directory",
			src:     files.FilePath{Dir: "/nonexistent/dir", Name: "file.txt"},
			dest:    files.FilePath{Dir: destDir, Name: "dest3.txt"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := files.MoveFile(tt.src, tt.dest)

			if (err != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify source file no longer exists
				testutil.AssertFileNotExists(t, tt.src.FullPath())

				// Verify destination file exists with correct content
				testutil.AssertFileExists(t, tt.dest.FullPath())
				testutil.AssertFileContent(t, tt.dest.FullPath(), testContent)
			}
		})
	}
}

// TestDeleteFile tests the DeleteFile function.
// It verifies file deletion operations work correctly in various scenarios.
func TestDeleteFile(t *testing.T) {
	// Create a temporary test directory
	testDir := testutil.CreateTestDirectory(t)

	// Create a test file
	testFile := "test.txt"
	err := testutil.CreateTestFile(t, testDir, testFile, "test content")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    files.FilePath
		wantErr bool
	}{
		{
			name: "existing file",
			path: files.FilePath{Dir: testDir, Name: testFile},
		},
		{
			name:    "non-existent file",
			path:    files.FilePath{Dir: testDir, Name: "nonexistent.txt"},
			wantErr: true,
		},
		{
			name:    "invalid directory",
			path:    files.FilePath{Dir: "/nonexistent/dir", Name: "file.txt"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := files.DeleteFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file no longer exists
				testutil.AssertFileNotExists(t, tt.path.FullPath())
			}
		})
	}
}
