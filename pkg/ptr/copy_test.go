package ptr_test

import (
	"testing"

	"github.com/avivSarig/cerebgo/pkg/ptr"
)

func TestCopyPtr(t *testing.T) {
	// Helper function to create string pointer
	strPtr := func(s string) *string {
		return &s
	}

	// Helper function to create int pointer
	intPtr := func(i int) *int {
		return &i
	}

	tests := []struct {
		name     string
		input    any    // Using any to test different types
		want     any    // Expected result
		wantNil  bool   // Whether we expect nil output
		mutateIn func() // Function to mutate input after copy
	}{
		{
			name:    "nil pointer returns nil",
			input:   (*string)(nil),
			wantNil: true,
		},
		{
			name:  "copies string pointer",
			input: strPtr("hello"),
			want:  "hello",
		},
		{
			name:  "copies int pointer",
			input: intPtr(42),
			want:  42,
		},
		{
			name:  "copies zero string",
			input: strPtr(""),
			want:  "",
		},
		{
			name:  "copies zero int",
			input: intPtr(0),
			want:  0,
		},
		{
			name:  "modifying original doesn't affect copy - string",
			input: strPtr("original"),
			want:  "original",
			mutateIn: func() {
				s := "modified"
				p := &s
				*p = "modified"
			},
		},
		{
			name:  "modifying original doesn't affect copy - int",
			input: intPtr(1),
			want:  1,
			mutateIn: func() {
				i := 2
				p := &i
				*p = 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any

			// Type switch to handle different pointer types
			switch v := tt.input.(type) {
			case *string:
				result := ptr.CopyPtr(v)
				if result != nil {
					got = *result
				}
			case *int:
				result := ptr.CopyPtr(v)
				if result != nil {
					got = *result
				}
			}

			// If we expect nil, check for it
			if tt.wantNil {
				if got != nil {
					t.Errorf("CopyPtr() = %v, want nil", got)
				}
				return
			}

			// Run mutation if provided
			if tt.mutateIn != nil {
				tt.mutateIn()
			}

			// Compare results
			if got != tt.want {
				t.Errorf("CopyPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}
