package tasks_test

import (
	"testing"

	"github.com/avivSarig/cerebgo/pkg/tasks"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

// Test configuration constants.
const validConfig = `
paths:
    base:
        tasks: /tasks
        journal: /journal
    subdirs:
        tasks:
            completed: ${paths.base.tasks}/completed
settings:
    patterns:
        date_format: "YYYY-MM-DD"
`

const invalidConfig = `
paths: [invalid: yaml
`

func TestInitialization(t *testing.T) {
	testutil.SetEnv(t, "DATA_PATH", "/test/data")

	tests := []struct {
		name       string
		configYAML string
		setupEnv   func(t *testing.T)
		wantErr    bool
	}{
		{
			name:       "valid configuration loads successfully",
			configYAML: validConfig,
			setupEnv: func(t *testing.T) {
				dir := testutil.SetupConfigDir(t, validConfig)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: false,
		},
		{
			name:       "missing configuration file returns error",
			configYAML: "",
			setupEnv: func(t *testing.T) {
				dir := testutil.CreateTestDirectory(t)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: true,
		},
		{
			name:       "invalid YAML returns error",
			configYAML: invalidConfig,
			setupEnv: func(t *testing.T) {
				dir := testutil.SetupConfigDir(t, invalidConfig)
				testutil.SetConfigPath(t, dir)
			},
			wantErr: true,
		},
		{
			name:       "missing required configuration fields returns error",
			configYAML: "foo: bar",
			setupEnv: func(t *testing.T) {
				dir := testutil.SetupConfigDir(t, "foo: bar")
				testutil.SetConfigPath(t, dir)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state before each test
			tasks.ResetForTesting()

			// Setup test environment
			if tt.setupEnv != nil {
				tt.setupEnv(t)
			}

			// Test initialization
			err := tasks.Initialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Test GetConfig
				config, err := tasks.GetConfig()
				if err != nil {
					t.Errorf("GetConfig() unexpected error = %v", err)
					return
				}
				if config == nil {
					t.Error("GetConfig() returned nil config after successful initialization")
					return
				}

				// Verify required configuration fields
				if !config.IsSet("paths.base.tasks") {
					t.Error("required configuration field 'paths.base.tasks' is missing")
				}
				if !config.IsSet("settings.patterns.date_format") {
					t.Error("required configuration field 'settings.patterns.date_format' is missing")
				}
			}
		})
	}
}
