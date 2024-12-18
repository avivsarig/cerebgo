package testutil

import "fmt"

type ValidationResult struct {
	Field   string
	IsValid bool
	Message string
	Got     any
	Want    any
}

func (v ValidationResult) ToString() string {
	if !v.IsValid {
		return fmt.Sprintf("%s: got=%v, want=%v. %s",
			v.Field, v.Got, v.Want, v.Message)
	}
	return ""
}

type Comparer[T any] func(T, T) bool

func GenerateErrorMessages(results []ValidationResult) []string {
	var messages []string
	for _, result := range results {
		if !result.IsValid {
			messages = append(messages, result.ToString())
		}
	}
	return messages
}
