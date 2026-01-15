package gcid

import "cmp"

type IdtOrdered interface {
	cmp.Ordered
	String() string
}

type Identifiable[T IdtOrdered] interface {
	Idt() T
}
