package tasks

import (
	"fmt"
	"sync"

	"github.com/avivSarig/cerebgo/config"
	"github.com/spf13/viper"
)

var (
	configuration *viper.Viper
	initOnce      sync.Once
	initErr       error
)

// validateConfiguration ensures all required configuration fields are present.
// Returns an error if any required field is missing.
func validateConfiguration(v *viper.Viper) error {
	requiredPaths := []string{
		"paths.base.tasks",
		"paths.base.journal",
		"settings.patterns.date_format",
	}

	for _, path := range requiredPaths {
		if !v.IsSet(path) {
			return fmt.Errorf("missing required configuration field: %s", path)
		}
	}

	return nil
}

// Initialize loads and validates the configuration if it hasn't been loaded yet.
// It is safe to call multiple times - only the first call will perform initialization.
func Initialize() error {
	initOnce.Do(func() {
		var v *viper.Viper
		v, initErr = config.LoadConfig()
		if initErr == nil {
			// Only validate if loading was successful
			initErr = validateConfiguration(v)
			if initErr == nil {
				configuration = v
			}
		}
	})
	return initErr
}

// GetConfig returns the current configuration.
// It will initialize if needed.
func GetConfig() (*viper.Viper, error) {
	if err := Initialize(); err != nil {
		return nil, err
	}
	return configuration, nil
}

// For testing purposes only.
func ResetForTesting() {
	configuration = nil
	initOnce = sync.Once{}
	initErr = nil
}
