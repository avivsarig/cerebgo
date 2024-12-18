package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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
	v.AddConfigPath("/etc/myapp")

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return v, nil
}
