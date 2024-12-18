package testutil

import "time"

func StringComparer(got, want string) bool  { return got == want }
func TimeComparer(got, want time.Time) bool { return got.Equal(want) }
func BoolComparer(got, want bool) bool      { return got == want }
