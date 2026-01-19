package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type QueryPaginatedRequest[I gcfield.IdtOrdered, E gcfield.Identifiable[I]] struct {
	QueryItems                  string
	ArgsQueryItems              []any
	ConverterQueryItems         func(row SqlRow) (entity E, err error)
	QueryFirstLastIdts          string
	ArgsQueryFirstLastIdts      []any
	ConverterQueryFirstLastIdts func(row SqlRow) (firstIdt, lastIdt gcopt.Optional[I], err error)
}

type PagingCriteria[I gcfield.IdtOrdered] struct {
	Idt        I
	IsValidIdt bool
	Filter     string
	OrderBy    string
}
