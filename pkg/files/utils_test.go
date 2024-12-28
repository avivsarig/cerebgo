package files_test

import (
	"path/filepath"
	"testing"

	"github.com/avivSarig/cerebgo/pkg/files"
)

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
