package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type NewPagingCriteriaParams struct {
	Idt   ColumnCriteria
	Field gcopt.Optional[ColumnCriteria]
}

type QueryPaginatedParams[I gcfield.IdtOrdered, E gcfield.Identifiable[I]] struct {
	QueryItems                  string
	ArgsQueryItems              []any
	ConverterQueryItems         func(row SqlRow) (entity E, err error)
	QueryFirstLastIdts          string
	ArgsQueryFirstLastIdts      []any
	ConverterQueryFirstLastIdts func(row SqlRow) (firstIdt, lastIdt gcopt.Optional[I], err error)
}
