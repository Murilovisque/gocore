package gcpag

import "cmp"

type IdtOrdered interface {
	cmp.Ordered
	String() string
	// IsValid() bool
}

type Identifiable[T IdtOrdered] interface {
	Idt() T
}

type Field interface {
	IsValid() bool
}
