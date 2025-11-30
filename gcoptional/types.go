package gcoptional

type Optional[T any] struct {
	vl     T
	exists bool
}

func (o Optional[T]) IsPresent() bool {
	return o.exists
}

func (o Optional[T]) Take() (T, bool) {
	return o.vl, o.exists
}

func (o Optional[T]) MustTake() T {
	if o.exists {
		return o.vl
	}
	panic("gcoptional: attempt to MustTake() value from invalid Optional")
}

func (o Optional[T]) TakeOrElse(orElseFunc func() T) T {
	if o.exists {
		return o.vl
	}
	return orElseFunc()
}

func (o Optional[T]) TakeOrError(orElseErrFunc func() error) (T, error) {
	if o.exists {
		return o.vl, nil
	}
	return o.vl, orElseErrFunc()
}

func FromPointer[T any](vl *T) Optional[T] {
	if vl == nil {
		return None[T]()
	}
	return FromValue(*vl)
}

func FromValue[T any](vl T) Optional[T] {
	return Optional[T]{
		vl:     vl,
		exists: true,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}
