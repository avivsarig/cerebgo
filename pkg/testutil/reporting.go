package testutil

import "testing"

func ReportResults(t *testing.T, results []ValidationResult) {
	t.Helper()
	for _, msg := range GenerateErrorMessages(results) {
		t.Error(msg)
	}
}
