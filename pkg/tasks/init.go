package tasks

import (
	"log"

	"github.com/avivSarig/cerebgo/config"
	"github.com/spf13/viper"
)

var configuration *viper.Viper

func init() {
	var err error
	configuration, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}
