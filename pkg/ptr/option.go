package ptr

// Option[T] represents an optional value that may or may not be present.
// Use Some(value) to create a valid Option, or None[T]() for an empty Option.
type Option[T any] struct {
	value T
	valid bool
}

// Some creates an Option containing a valid value of any type T.
// Option[T] represents an optional value that may or may not be present.
//
// Parameters:
//   - value: The value to wrap in the Option
//
// Returns:
//
//	Option[T] with valid=true and the provided value
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

// None creates an empty Option of type T.
// Used to represent absence of a value.
//
// Returns:
//
//	Option[T] with valid=false and zero value of T
func None[T any]() Option[T] {
	return Option[T]{valid: false}
}

// IsValid returns whether the Option contains a valid value.
//
// Returns:
//
//	true if Option contains a value, false if None
func (o Option[T]) IsValid() bool {
	return o.valid
}

// Value returns the contained value.
// Panics if Option is None - check IsValid() first.
//
// Returns:
//
//	The contained value of type T
func (o Option[T]) Value() T {
	if !o.valid {
		panic("attempted to access Value() of None option")
	}
	return o.value
}
