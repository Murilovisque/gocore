package gcopt

func OfPtr[T any](vl *T) Optional[T] {
	if vl == nil {
		return Empty[T]()
	}
	return Of(*vl)
}

func Of[T any](vl T) Optional[T] {
	return Optional[T]{
		value:  vl,
		exists: true,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{exists: false}
}

func Map[T any, U any](o Optional[T], f func(T) U) Optional[U] {
	if !o.exists {
		return Empty[U]()
	}
	return Of(f(o.value))
}
