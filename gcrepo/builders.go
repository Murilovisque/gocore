package gcrepo

import (
	"context"
	"fmt"
	"slices"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcpag"
)

func PlaceHolderRange(s SqlSyntax, startAt, stopBefore int) []any {
	phs := make([]any, 0, stopBefore-startAt)
	for i := startAt; i < stopBefore; i++ {
		phs = append(phs, s.PlaceHolder(i))
	}
	return phs
}

func NewPagingCriteria[I gcfield.IdtOrdered](pagReq gcpag.PaginatedRequest[I], params NewPagingCriteriaParams) PagingCriteria[I] {
	idtValue, ok := pagReq.Idt.Take()
	if !ok {
		return PagingCriteria[I]{
			IsValid: false,
			OrderBy: fmt.Sprintf("%s %s", params.Idt.Column, pagReq.Order.String()),
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
	if _, ok := pagReq.Field.Take(); ok && params.Field.IsPresent() {
		sqlFld := params.Field.MustTake()
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
			IsValid: true,
			Filter:  fmt.Sprintf("%s %s %s and %s %s %s", params.Idt.Column, cmpSignal, params.Idt.PlaceHolder, sqlFld.Column, cmpSignalField, sqlFld.PlaceHolder),
			OrderBy: fmt.Sprintf("%s %s, %s %s", params.Idt.Column, idtOrder.String(), sqlFld.Column, idtOrder.String()),
		}
	}
	return PagingCriteria[I]{
		Idt:     idtValue,
		IsValid: true,
		Filter:  fmt.Sprintf("%s %s %s", params.Idt.Column, cmpSignal, params.Idt.PlaceHolder),
		OrderBy: fmt.Sprintf("%s %s", params.Idt.Column, idtOrder.String()),
	}
}

func QueryPaginated[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](ctx context.Context, executor SqlExecutor, pagReq gcpag.PaginatedRequest[I], req QueryPaginatedParams[I, E]) (gcpag.PaginatedResponse[I, E], error) {
	var res gcpag.PaginatedResponse[I, E]
	rowsItens, err := executor.Query(ctx, req.QueryItems, req.ArgsQueryItems...)
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
	rowFirstLast := executor.QueryRow(ctx, req.QueryFirstLastIdts, req.ArgsQueryFirstLastIdts...)
	minIdt, maxIdt, err := req.ConverterQueryFirstLastIdts(rowFirstLast)
	if err != nil {
		return res, fmt.Errorf("gcrepo: convert first last idt failed. Cause %w", err)
	}
	return gcpag.NewResponse(pagReq, items, minIdt, maxIdt), nil
}
