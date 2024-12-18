package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	destDir := filepath.Join(tmpDir, "dest")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := "test.txt"
	srcPath := filepath.Join(srcDir, testFile)
	if err := os.WriteFile(srcPath, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		src     string
		dest    string
		file    string
		wantErr bool
	}{
		{
			name:    "successful move",
			src:     srcDir,
			dest:    destDir,
			file:    testFile,
			wantErr: false,
		},
		{
			name:    "file not found",
			src:     srcDir,
			dest:    destDir,
			file:    "nonexistent.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MoveFile(tt.src, tt.dest, tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
