package nullable

type Type[T any] struct {
	ptr   *T
	Valid bool
}

func TypeFrom[T any](v T) Type[T] {
	return NewType(&v, true)
}

func TypeFromPtr[T any](v *T) Type[T] {
	if v == nil {
		var zeroValue T
		return NewType(&zeroValue, false)
	}
	return NewType(v, true)
}

func NewType[T any](v *T, valid bool) Type[T] {
	return Type[T]{
		ptr:   v,
		Valid: valid,
	}
}

func (s Type[T]) Ptr() *T {
	return s.ptr
}

func (s Type[T]) Value() T {
	return *s.ptr
}

func (s Type[T]) IsZero() bool {
	return !s.Valid
}
