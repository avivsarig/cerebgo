package main

import (
	"log"
	"time"

	"github.com/avivSarig/cerebgo/config"
	"github.com/avivSarig/cerebgo/pkg/tasks"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get current time for task processing
	now := time.Now()

	// Process all tasks
	if err := tasks.ProcessAllTasks(now, cfg); err != nil {
		log.Fatalf("Failed to process tasks: %v", err)
	}
}
