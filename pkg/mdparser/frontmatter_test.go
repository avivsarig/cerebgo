package mdparser_test

import (
	"testing"
	"time"

	"github.com/avivSarig/cerebgo/pkg/mdparser"
	"github.com/avivSarig/cerebgo/pkg/ptr"
	"github.com/avivSarig/cerebgo/pkg/testutil"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		fm       mdparser.Frontmatter
		key      string
		want     string
		wantBool bool
	}{
		{
			name:     "valid string",
			fm:       mdparser.Frontmatter{"title": "Hello"},
			key:      "title",
			want:     "Hello",
			wantBool: true,
		},
		{
			name:     "empty string",
			fm:       mdparser.Frontmatter{"title": ""},
			key:      "title",
			want:     "",
			wantBool: true,
		},
		{
			name:     "unicode string",
			fm:       mdparser.Frontmatter{"title": "שלום"},
			key:      "title",
			want:     "שלום",
			wantBool: true,
		},
		{
			name:     "key doesn't exist",
			fm:       mdparser.Frontmatter{},
			key:      "title",
			want:     "",
			wantBool: false,
		},
		{
			name:     "wrong type",
			fm:       mdparser.Frontmatter{"title": true},
			key:      "title",
			want:     "",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := mdparser.GetString(tt.fm, tt.key)
			results := []testutil.ValidationResult{
				testutil.ValidateEqual("value", got, tt.want),
				testutil.ValidateEqual("ok", ok, tt.wantBool),
			}
			testutil.ReportResults(t, results)
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		fm       mdparser.Frontmatter
		key      string
		want     bool
		wantBool bool
	}{
		{
			name:     "valid true",
			fm:       mdparser.Frontmatter{"draft": true},
			key:      "draft",
			want:     true,
			wantBool: true,
		},
		{
			name:     "valid false",
			fm:       mdparser.Frontmatter{"draft": false},
			key:      "draft",
			want:     false,
			wantBool: true,
		},
		{
			name:     "key doesn't exist",
			fm:       mdparser.Frontmatter{},
			key:      "draft",
			want:     false,
			wantBool: false,
		},
		{
			name:     "wrong type",
			fm:       mdparser.Frontmatter{"draft": "true"},
			key:      "draft",
			want:     false,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := mdparser.GetBool(tt.fm, tt.key)
			results := []testutil.ValidationResult{
				testutil.ValidateEqual("value", got, tt.want),
				testutil.ValidateEqual("ok", ok, tt.wantBool),
			}
			testutil.ReportResults(t, results)
		})
	}
}

func TestGetTime(t *testing.T) {
	validTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	validTimeStr := validTime.Format(time.RFC3339)

	tests := []struct {
		name     string
		fm       mdparser.Frontmatter
		key      string
		want     ptr.Option[time.Time]
		wantBool bool
	}{
		{
			name:     "valid time",
			fm:       mdparser.Frontmatter{"created": validTimeStr},
			key:      "created",
			want:     ptr.Some(validTime),
			wantBool: true,
		},
		{
			name:     "invalid format",
			fm:       mdparser.Frontmatter{"created": "2024-01-01"},
			key:      "created",
			want:     ptr.None[time.Time](),
			wantBool: false,
		},
		{
			name:     "key doesn't exist",
			fm:       mdparser.Frontmatter{},
			key:      "created",
			want:     ptr.None[time.Time](),
			wantBool: false,
		},
		{
			name:     "wrong type",
			fm:       mdparser.Frontmatter{"created": true},
			key:      "created",
			want:     ptr.None[time.Time](),
			wantBool: false,
		},
		{
			name:     "extreme time",
			fm:       mdparser.Frontmatter{"created": "9999-12-31T23:59:59Z"},
			key:      "created",
			want:     ptr.Some(time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)),
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, ok := mdparser.GetTime(tt.fm, tt.key)
			got := ptr.None[time.Time]()
			if ok {
				got = ptr.Some(gotTime)
			}

			results := []testutil.ValidationResult{
				testutil.ValidateOptional(
					"time",
					got,
					tt.want,
					testutil.TimeComparer,
				),
			}
			testutil.ReportResults(t, results)
		})
	}
}
