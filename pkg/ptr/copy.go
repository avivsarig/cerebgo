package ptr

func CopyPtr[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	copied := *ptr
	return &copied
}
