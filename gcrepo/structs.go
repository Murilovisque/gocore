package gcrepo

import (
	"github.com/Murilovisque/gocore/gcfield"
)

type PagingCriteria[I gcfield.IdtOrdered] struct {
	Idt   I
	Query string
	Args  []any
	// Condition string //TODO: remove
	// OrderBy string
	// Limit   string
	// IsSortedByField bool   //TODO: remove
	// OrderByFirst    string //TODO: rename TO OrderByFirst, bellow too
	// QueryLast       string
}

type ColumnCriteria struct {
	Column      string //T0DO: rename
	PlaceHolder int
}

type SubQueryPaginatedFieldOrdered struct { //TOD: maybe rename
	SubQuery   string
	ColumnName string
}

type sqlRowWrapper struct {
	SqlRow
	fnWrapper func(args ...any) ([]any, error)
}

func (s *sqlRowWrapper) Scan(args ...any) error {
	args, err := s.fnWrapper(args)
	if err != nil {
		return err
	}
	return s.SqlRow.Scan(args)
}
