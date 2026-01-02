package gcpag

import "cmp"

type IdtOrdered interface {
	cmp.Ordered
	String() string
	// IsValid() bool
}

type ModelOrderable[T IdtOrdered] interface {
	OrderableIdt() T
}

type Field interface {
	IsValid() bool
}
