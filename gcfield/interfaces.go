package gcfield

import (
	"cmp"

	"github.com/Murilovisque/gocore/gcopt"
)

type IdtOrdered interface {
	cmp.Ordered
	String() string
}

type Identifiable[T IdtOrdered] interface {
	Idt() T
}

type FieldNameOrderedParser interface {
	ParseFieldNameOrdered(name string) (gcopt.Optional[FieldNameOrdered], error)
}
