package gcrepo

import (
	"context"
	"fmt"
	"slices"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
	"github.com/Murilovisque/gocore/gcpag"
)

func PlaceHolderRange(s SqlSyntax, startAt, stopBefore int) []any {
	phs := make([]any, 0, stopBefore-startAt)
	for i := startAt; i < stopBefore; i++ {
		phs = append(phs, s.PlaceHolder(i))
	}
	return phs
}

func newPagingCriteria[I gcfield.IdtOrdered](syntax SqlSyntax, pagReq gcpag.PaginatedRequest[I], params NewPagingCriteriaParams) PagingCriteria[I] {
	idtValue, ok := pagReq.Idt.Take()
	if !ok {
		return PagingCriteria[I]{
			OrderBy: fmt.Sprintf("order by %s %s", params.Idt.Column, pagReq.Order.String()),
			Limit:   syntax.LimitStatement(params.Idt.PlaceHolder),
			Args:    []any{pagReq.Size},
		}
	}
	var cmpSignal string
	var idtOrder gcpag.Order
	switch pagReq.Orientation {
	case gcpag.NextPage:
		idtOrder = gcpag.Asc
		switch pagReq.StartPosition {
		case gcpag.StartAt:
			cmpSignal = ">="
		case gcpag.AfterAt:
			cmpSignal = ">"
		}
	case gcpag.PreviousPage:
		idtOrder = gcpag.Desc
		switch pagReq.StartPosition {
		case gcpag.StartAt:
			cmpSignal = "<="
		case gcpag.AfterAt:
			cmpSignal = "<"
		}
	}
	if pagReq.Order == gcpag.Desc {
		switch cmpSignal {
		case ">":
			cmpSignal = "<"
		case ">=":
			cmpSignal = "<="
		case "<":
			cmpSignal = ">"
		case "<=":
			cmpSignal = ">="
		}
		idtOrder = idtOrder.Reverse()
	}
	if params.Field.IsPresent() && pagReq.Field.IsPresent() {
		sqlFld := params.Field.MustTake()
		paramFld := pagReq.Field.MustTake()
		var cmpSignalField string
		switch pagReq.Orientation {
		case gcpag.NextPage:
			cmpSignalField = ">="
		case gcpag.PreviousPage:
			cmpSignalField = "<="
		}
		return PagingCriteria[I]{
			Idt:     idtValue,
			Field:   pagReq.Field,
			Where:   fmt.Sprintf("%s %s %s and %s %s %s", params.Idt.Column, cmpSignal, syntax.PlaceHolder(params.Idt.PlaceHolder), sqlFld.Column, cmpSignalField, syntax.PlaceHolder(sqlFld.PlaceHolder)),
			OrderBy: fmt.Sprintf("order by %s %s, %s %s", params.Idt.Column, idtOrder.String(), sqlFld.Column, idtOrder.String()),
			Limit:   syntax.LimitStatement(sqlFld.PlaceHolder + 1),
			Args:    []any{idtValue, paramFld.Value(), pagReq.Size},
		}
	}
	return PagingCriteria[I]{
		Idt:     idtValue,
		Where:   fmt.Sprintf("%s %s %s", params.Idt.Column, cmpSignal, syntax.PlaceHolder(params.Idt.PlaceHolder)),
		OrderBy: fmt.Sprintf("order by %s %s", params.Idt.Column, idtOrder.String()),
		Limit:   syntax.LimitStatement(params.Idt.PlaceHolder + 1),
		Args:    []any{idtValue, pagReq.Size},
	}
}

func QueryPaginated[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](ctx context.Context, executor SqlExecutor, pagReq gcpag.PaginatedRequest[I], req QueryPaginatedParams[I, E]) (gcpag.PaginatedResponse[I, E], error) {
	var res gcpag.PaginatedResponse[I, E]
	var query string
	syntax := executor.Syntax()
	pagCrit := newPagingCriteria(syntax, pagReq, NewPagingCriteriaParams{
		Idt: ColumnCriteria{Column: req.IdtColumn, PlaceHolder: req.LastQueryPlaceHolder + 1},
		Field: gcopt.Map(req.FieldColumn, func(fldColumn string) ColumnCriteria {
			return ColumnCriteria{
				Column:      fldColumn,
				PlaceHolder: req.LastQueryPlaceHolder + 2,
			}
		}),
	})
	if pagCrit.Where != "" {
		query = fmt.Sprintf("%s and %s %s %s", req.QueryItems, pagCrit.Where, pagCrit.OrderBy, pagCrit.Limit)
	} else {
		query = fmt.Sprintf("%s %s %s", req.QueryItems, pagCrit.OrderBy, pagCrit.Limit)
	}
	args := append(req.QueryArgs, pagCrit.Args...)
	rowsItens, err := executor.Query(ctx, query, args...)
	if err != nil {
		return res, fmt.Errorf("gcrepo: query itens paginated failed. Cause %w", err)
	}
	defer rowsItens.Close()
	var items []E
	for rowsItens.Next() {
		m, err := req.ConverterQueryItems(rowsItens)
		if err != nil {
			return res, fmt.Errorf("gcrepo: convert item of the rows failed. Cause %w", err)
		}
		items = append(items, m)
	}
	if err = rowsItens.Err(); err != nil {
		return res, fmt.Errorf("gcrepo: convert item after rows iteration failed. Cause %w", err)
	}
	if pagReq.Orientation == gcpag.PreviousPage {
		slices.Reverse(items)
	}
	rowFirstLast := executor.QueryRow(ctx, req.QueryFirstLastIdts, req.QueryArgs...)
	minIdt, maxIdt, err := req.ConverterQueryFirstLastIdts(rowFirstLast)
	if err != nil {
		return res, fmt.Errorf("gcrepo: convert first last idt failed. Cause %w", err)
	}
	return gcpag.NewResponse(pagReq, items, minIdt, maxIdt), nil
}
