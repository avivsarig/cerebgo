package ptr

type Option[T any] struct {
	value T
	valid bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

func None[T any]() Option[T] {
	return Option[T]{valid: false}
}

func (o Option[T]) IsValid() bool {
	return o.valid
}

func (o Option[T]) Value() T {
	return o.value
}
