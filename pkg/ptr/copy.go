package ptr

// CopyPtr creates a deep copy of a pointer to T.
// Useful for avoiding shared pointer references.
//
// Parameters:
//   - ptr: Source pointer to copy
//
// Returns:
//
//	New pointer to a copy of the value, or nil if ptr is nil
func CopyPtr[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	copied := *ptr
	return &copied
}
