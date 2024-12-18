package testutil

import "github.com/avivSarig/cerebgo/pkg/ptr"

func ValidateOptional[T any](
	field string,
	got, want ptr.Option[T],
	comparer Comparer[T],
) ValidationResult {
	if got.IsValid() != want.IsValid() {
		return CreateValidationError(
			field,
			got.IsValid(),
			want.IsValid(),
			"validity mismatch",
		)
	}

	if got.IsValid() && !comparer(got.Value(), want.Value()) {
		return CreateValidationError(
			field,
			got.Value(),
			want.Value(),
			"value mismatch",
		)
	}

	return CreateValidSuccess(field)
}

func ValidateEqual[T comparable](
	field string,
	got, want T,
) ValidationResult {
	if got != want {
		return CreateValidationError(field, got, want, "values not equal")
	}
	return CreateValidSuccess(field)
}
