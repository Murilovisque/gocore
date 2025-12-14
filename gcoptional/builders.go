package gcoptional

func FromPointer[T any](vl *T) Optional[T] {
	if vl == nil {
		return EmtpyValue[T]()
	}
	return FromValue(*vl)
}

func FromValue[T any](vl T) Optional[T] {
	return Optional[T]{
		vl:     vl,
		exists: true,
	}
}

func EmtpyValue[T any]() Optional[T] {
	return Optional[T]{}
}
