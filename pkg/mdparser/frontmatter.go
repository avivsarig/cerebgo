package mdparser

import (
	"time"
)

// Frontmatter represents markdown document metadata as key-value pairs.
type Frontmatter map[string]interface{}

// GetString extracts a string value from Frontmatter.
//
// Parameters:
//   - fm: Frontmatter to extract from
//   - key: Key to look up
//
// Returns:
//
//	String value and true if exists and is string type, empty and false otherwise
func GetString(fm Frontmatter, key string) (string, bool) {
	val, ok := fm[key].(string)
	return val, ok
}

// GetBool extracts a boolean value from Frontmatter.
//
// Parameters:
//   - fm: Frontmatter to extract from
//   - key: Key to look up
//
// Returns:
//
//	Bool value and true if exists and is bool type, false and false otherwise
func GetBool(fm Frontmatter, key string) (bool, bool) {
	val, ok := fm[key].(bool)
	return val, ok
}

// GetTime extracts and parses a time value from Frontmatter.
// Expects RFC3339 formatted string.
//
// Parameters:
//   - fm: Frontmatter to extract from
//   - key: Key to look up
//
// Returns:
//
//	Parsed time.Time and true if value exists and is valid RFC3339, zero time and false otherwise
func GetTime(fm Frontmatter, key string) (time.Time, bool) {
	strVal, ok := fm[key].(string)
	if !ok {
		return time.Time{}, false
	}
	val, err := time.Parse(time.RFC3339, strVal)
	return val, err == nil
}
