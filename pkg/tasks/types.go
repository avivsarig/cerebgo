package tasks

import (
	"github.com/avivSarig/cerebgo/internal/models"
)

type ActionType int

const (
	ActionRetain ActionType = iota
	ActionArchive
	ActionUpdate
	ActionComplete
)

type TaskAction struct {
	Task    models.Task
	Action  ActionType
	Updates []TaskModifier
}

// Map is a generic higher-order function for transforming slices.
func Map[T, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}
