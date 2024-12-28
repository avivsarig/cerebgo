package ptr_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestOption(t *testing.T) {
	// Define a test struct to represent a custom type
	type customStruct struct {
		x int
		y string
	}

	timeNow := time.Now()
	structVal := customStruct{x: 1, y: "test"}

	tests := []struct {
		name      string
		setup     func() ptr.Option[any]   // Function to create the option we're testing
		wantValid bool                     // Whether we expect IsValid() to return true
		wantValue any                      // Expected value (ignored if wantValid is false)
		wantPanic bool                     // Whether we expect Value() to panic
		compare   func(got, want any) bool // Optional comparison function for Value()
	}{
		// Some basic cases
		{
			name:      "Some with string",
			setup:     func() ptr.Option[any] { return ptr.Some[any]("hello") },
			wantValid: true,
			wantValue: "hello",
		},
		{
			name:      "Some with empty string",
			setup:     func() ptr.Option[any] { return ptr.Some[any]("") },
			wantValid: true,
			wantValue: "",
		},
		{
			name:      "Some with zero int",
			setup:     func() ptr.Option[any] { return ptr.Some[any](0) },
			wantValid: true,
			wantValue: 0,
		},
		{
			name:      "Some with negative int",
			setup:     func() ptr.Option[any] { return ptr.Some[any](-42) },
			wantValid: true,
			wantValue: -42,
		},

		// Complex types
		{
			name:      "Some with struct",
			setup:     func() ptr.Option[any] { return ptr.Some[any](structVal) },
			wantValid: true,
			wantValue: structVal,
		},
		{
			name:      "Some with time.Time",
			setup:     func() ptr.Option[any] { return ptr.Some[any](timeNow) },
			wantValid: true,
			wantValue: timeNow,
		},
		{
			name:      "Some with slice",
			setup:     func() ptr.Option[any] { return ptr.Some[any]([]int{1, 2, 3}) },
			wantValid: true,
			wantValue: []int{1, 2, 3},
			compare: func(got, want any) bool {
				g, ok1 := got.([]int)
				w, ok2 := want.([]int)
				if !ok1 || !ok2 {
					return false
				}
				return reflect.DeepEqual(g, w)
			},
		},

		// None cases
		{
			name:      "None string",
			setup:     func() ptr.Option[any] { return ptr.None[any]() },
			wantValid: false,
			wantPanic: true,
		},
		{
			name:      "Zero value option",
			setup:     func() ptr.Option[any] { var opt ptr.Option[any]; return opt },
			wantValid: false,
			wantPanic: true,
		},

		// Edge cases
		{
			name:      "Some with nil interface",
			setup:     func() ptr.Option[any] { return ptr.Some[any](nil) },
			wantValid: true,
			wantValue: nil,
		},
		{
			name: "Some with pointer",
			setup: func() ptr.Option[any] {
				testInt := 42
				return ptr.Some[any](&testInt)
			},
			wantValid: true,
			wantValue: func() any {
				testInt := 42
				return &testInt
			}(),
			compare: func(got, want any) bool {
				g, ok1 := got.(*int)
				w, ok2 := want.(*int)
				if !ok1 || !ok2 {
					return false
				}
				return *g == *w
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.setup()

			if got := opt.IsValid(); got != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValid)
			}

			if tt.wantPanic {
				testutil.AssertPanics(t, func() { opt.Value() })
			} else if tt.wantValid {
				if tt.compare != nil {
					if !tt.compare(opt.Value(), tt.wantValue) {
						t.Errorf("Value() = %v, want %v", opt.Value(), tt.wantValue)
					}
				} else if got := opt.Value(); got != tt.wantValue {
					t.Errorf("Value() = %v, want %v", got, tt.wantValue)
				}
			}
		})
	}
}
