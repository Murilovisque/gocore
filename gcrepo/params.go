package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type NewPagingCriteriaParams struct { //TODO: remove
	Idt   ColumnCriteria
	Field gcopt.Optional[ColumnCriteria] //TODO: rename
}

type QueryPaginatedParams[I gcfield.IdtOrdered, E gcfield.Identifiable[I]] struct {
	QueryItems          string
	QueryArgs           []any
	ConverterQueryItems func(row SqlRow) (entity E, err error)
	// QueryFirstLastIdts          string //TODO: remove
	ConverterQueryFirstLastIdts func(row SqlRow) (firstIdt, lastIdt gcopt.Optional[I], err error) //TODO: remove
	IdtColumn                   string
	FieldColumn                 gcopt.Optional[func(fld gcfield.FieldNameOrdered, placeHolder int) gcopt.Optional[SubQueryPaginatedFieldOrdered]] //TODO: rename, chnage placeHolder int to string already adapted
	LastQueryPlaceHolder        int
}
