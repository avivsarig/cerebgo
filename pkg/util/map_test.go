package util_test

import (
	"strings"
	"testing"

	"github.com/avivSarig/cerebgo/pkg/util"
)

func TestMap(t *testing.T) {
	// Test case struct that can handle any input/output types
	type testCase[T, U any] struct {
		name     string
		input    []T
		f        func(T) U
		expected []U
	}

	// Run sub-tests for different types and scenarios
	t.Run("integers", func(t *testing.T) {
		cases := []testCase[int, int]{
			{
				name:     "empty slice",
				input:    []int{},
				f:        func(x int) int { return x * 2 },
				expected: []int{},
			},
			{
				name:     "double numbers",
				input:    []int{1, 2, 3, 4, 5},
				f:        func(x int) int { return x * 2 },
				expected: []int{2, 4, 6, 8, 10},
			},
			{
				name:     "square numbers",
				input:    []int{1, 2, 3, 4},
				f:        func(x int) int { return x * x },
				expected: []int{1, 4, 9, 16},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := util.Map(tc.input, tc.f)
				if len(got) != len(tc.expected) {
					t.Errorf("length mismatch: got %v, want %v", len(got), len(tc.expected))
					return
				}
				for i := range got {
					if got[i] != tc.expected[i] {
						t.Errorf("at index %d: got %v, want %v", i, got[i], tc.expected[i])
					}
				}
			})
		}
	})

	t.Run("string transformations", func(t *testing.T) {
		cases := []testCase[string, string]{
			{
				name:     "empty strings",
				input:    []string{},
				f:        strings.ToUpper,
				expected: []string{},
			},
			{
				name:     "to upper",
				input:    []string{"hello", "world"},
				f:        strings.ToUpper,
				expected: []string{"HELLO", "WORLD"},
			},
			{
				name:     "append suffix",
				input:    []string{"test", "case"},
				f:        func(s string) string { return s + "_suffix" },
				expected: []string{"test_suffix", "case_suffix"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := util.Map(tc.input, tc.f)
				if len(got) != len(tc.expected) {
					t.Errorf("length mismatch: got %v, want %v", len(got), len(tc.expected))
					return
				}
				for i := range got {
					if got[i] != tc.expected[i] {
						t.Errorf("at index %d: got %v, want %v", i, got[i], tc.expected[i])
					}
				}
			})
		}
	})

	t.Run("type conversion", func(t *testing.T) {
		cases := []testCase[int, string]{
			{
				name:     "int to string",
				input:    []int{1, 2, 3},
				f:        func(x int) string { return strings.Repeat("*", x) },
				expected: []string{"*", "**", "***"},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := util.Map(tc.input, tc.f)
				if len(got) != len(tc.expected) {
					t.Errorf("length mismatch: got %v, want %v", len(got), len(tc.expected))
					return
				}
				for i := range got {
					if got[i] != tc.expected[i] {
						t.Errorf("at index %d: got %v, want %v", i, got[i], tc.expected[i])
					}
				}
			})
		}
	})
}
