package nullable

type Type[T any] struct {
	Value T
	Valid bool
}

func TypeFrom[T any](v T) Type[T] {
	return NewType(v, true)
}

func TypeFromPtr[T any](v *T) Type[T] {
	if v == nil {
		return NewType(*v, false)
	}
	return NewType(*v, true)
}

func NewType[T any](v T, valid bool) Type[T] {
	return Type[T]{
		Value: v,
		Valid: valid,
	}
}

func (s Type[T]) Ptr() *T {
	if !s.Valid {
		return nil
	}
	return &s.Value
}

func (s Type[T]) IsZero() bool {
	return !s.Valid
}
