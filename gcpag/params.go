package gcpag

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type ParseRequestFromHttpParams[I gcfield.IdtOrdered] struct {
	DefaultOrder Order
	DefaultSize  int
	IdtParser    func(string) (I, bool, error)
	Field        gcopt.Optional[*gcfield.FieldParser]
}
