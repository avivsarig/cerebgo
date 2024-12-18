package testutil

func CreateValidationError(field string, got, want any, message string) ValidationResult {
	return ValidationResult{
		Field:   field,
		IsValid: false,
		Message: message,
		Got:     got,
		Want:    want,
	}
}

func CreateValidSuccess(field string) ValidationResult {
	return ValidationResult{
		Field:   field,
		IsValid: true,
	}
}

func CombineResults(results ...ValidationResult) []ValidationResult {
	combined := make([]ValidationResult, len(results))
	copy(combined, results)
	return combined
}
