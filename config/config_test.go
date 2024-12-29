package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/avivSarig/cerebgo/config"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// Test Configuration constants.
const validConfig = `
key: value
nested:
  setting: 42
`

const invalidConfig = `
key: [invalid yaml`

// TestLoadConfig verifies LoadConfig functionality across different scenarios:
// - Loading from CONFIG_PATH environment variable
// - Loading from executable directory
// - Handling missing configuration
// - Handling invalid YAML.
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr bool
	}{
		{
			name: "loads from CONFIG_PATH",
			setup: func(t *testing.T) {
				dir := testutil.SetupConfigDir(t, validConfig)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: false,
		},
		{
			name: "loads from executable directory",
			setup: func(t *testing.T) {
				exe, err := os.Executable()
				if err != nil {
					t.Fatal(err)
				}
				configDir := filepath.Join(filepath.Dir(exe), "config")
				if err := os.MkdirAll(configDir, 0755); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					os.RemoveAll(configDir)
				})

				err = testutil.CreateTestFile(t, configDir, "config.yaml", validConfig)
				if err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
		},
		{
			name: "fails with missing config",
			setup: func(t *testing.T) {
				dir := testutil.CreateTestDirectory(t)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: true,
		},
		{
			name: "fails with invalid YAML",
			setup: func(t *testing.T) {
				dir := testutil.SetupConfigDir(t, invalidConfig)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			got, err := config.LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				testutil.ValidateConfig(t, got)
			}
		})
	}
}
