package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type PagingCriteria[I gcfield.IdtOrdered] struct {
	Idt     I
	Field   gcopt.Optional[gcfield.FieldParser]
	IsValid bool
	Filter  string
	OrderBy string
}

type ColumnCriteria struct {
	Column      string
	PlaceHolder string
}
