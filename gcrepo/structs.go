package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type PagingCriteria[I gcfield.IdtOrdered] struct {
	Idt     I
	Field   gcopt.Optional[gcfield.FieldParser]
	Where   string
	OrderBy string
	Limit   string
	Args    []any
}

type ColumnCriteria struct {
	Column      string
	PlaceHolder int
}
