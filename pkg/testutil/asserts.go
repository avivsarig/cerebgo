package testutil

import "testing"

func AssertPanics(t *testing.T, f func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected function to panic: %v", msgAndArgs...)
		}
	}()
	f()
}

// Optional: Add a version that checks panic message.
func AssertPanicsWithMessage(t *testing.T, f func(), expectedMsg string) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected function to panic with message: %s", expectedMsg)
			return
		} else if msg, ok := r.(string); !ok || msg != expectedMsg {
			t.Errorf("got panic message %v, want %s", r, expectedMsg)
		}
	}()
	f()
}
