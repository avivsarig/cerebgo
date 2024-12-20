package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig initializes and loads the application configuration using Viper.
// It searches for a "config.yaml" file in the following locations (in order):
//  1. Path specified by the CONFIG_PATH environment variable.
//  2. A "config" directory near the executable's location.
//  3. A "config" directory one level above the executable's location.
//
// Returns:
//   - *viper.Viper: Configured Viper instance on success.
//   - error: If the configuration file is missing or invalid.
func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		v.AddConfigPath(configPath)
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("couldn't determine executable path: %w", err)
	}
	exeDir := filepath.Dir(exe)

	v.AddConfigPath(filepath.Join(exeDir, "config"))
	v.AddConfigPath(filepath.Join(exeDir, "..", "config"))

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return v, nil
}
