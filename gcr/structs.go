package gcr

import (
	"github.com/Murilovisque/gocore/gcid"
	"github.com/Murilovisque/gocore/gcopt"
)

type QueryPaginatedRequest[I gcid.IdtOrdered, E gcid.Identifiable[I]] struct {
	QueryItems                  string
	ArgsQueryItems              []any
	ConverterQueryItems         func(row SqlRow) (entity E, err error)
	QueryFirstLastIdts          string
	ArgsQueryFirstLastIdts      []any
	ConverterQueryFirstLastIdts func(row SqlRow) (firstIdt, lastIdt gcopt.Optional[I], err error)
}

type PagingCriteria[I gcid.IdtOrdered] struct {
	Idt        I
	IsValidIdt bool
	Filter     string
	OrderBy    string
}
