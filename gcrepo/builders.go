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

func newPagingCriteria[I gcfield.IdtOrdered](syntax SqlSyntax, pagReq gcpag.PaginatedRequest[I], query string, idtColumn ColumnCriteria, sortFieldFn gcopt.Optional[func(fld gcfield.FieldNameOrdered, placeHolder int) gcopt.Optional[SubQueryPaginatedFieldOrdered]]) PagingCriteria[I] {
	idtValue, ok := pagReq.Idt.Take()
	sizeMoreNextElement := pagReq.Size + 1
	if !ok {
		// Without page parameter
		if sortFieldFn.IsPresent() && pagReq.SortField.IsPresent() {
			// Without page parameter and with sorted field
			sqlFld := sortFieldFn.MustTake() //rename here and bellow and more bellow
			paramFld := pagReq.SortField.MustTake()
			if sortFieldCriteria, ok := sqlFld(paramFld, idtColumn.PlaceHolder).Take(); ok {
				return PagingCriteria[I]{
					Query: fmt.Sprintf(
						`select *, false as is_previous_page from (%s) as sub_gorepo_pag
						order by sub_gorepo_pag.%s %s, sub_gorepo_pag.%s %s
						%s`,
						query,
						sortFieldCriteria.ColumnName, pagReq.Order.String(), idtColumn.Column, pagReq.Order.String(),
						syntax.LimitStatement(idtColumn.PlaceHolder)),
					Args: []any{sizeMoreNextElement},
				}
			}
		}
		return PagingCriteria[I]{
			Query: fmt.Sprintf(
				`select *, false as is_previous_page from (%s) as sub_gorepo_pag
				order by sub_gorepo_pag.%s %s
				%s`,
				query,
				idtColumn.Column, pagReq.Order.String(),
				syntax.LimitStatement(idtColumn.PlaceHolder)),
			Args: []any{sizeMoreNextElement},
		}
	}
	if sortFieldFn.IsPresent() && pagReq.SortField.IsPresent() {
		// With page parameter and sorted field
		sqlFld := sortFieldFn.MustTake()
		paramFld := pagReq.SortField.MustTake()
		if sortFieldCriteria, ok := sqlFld(paramFld, idtColumn.PlaceHolder).Take(); ok {
			var (
				cmpSignalTuple     CmpSign
				cmpSignalPrevious  CmpSign
				tupleOrder         gcpag.Order
				tupleOrderPrevious gcpag.Order
			)
			switch pagReq.Orientation {
			case gcpag.NextPage:
				tupleOrder = gcpag.Asc
				tupleOrderPrevious = gcpag.Desc
				cmpSignalTuple = CmpSignalGreater
				cmpSignalPrevious = CmpSignalLessOrEqual
			case gcpag.PreviousPage:
				tupleOrder = gcpag.Desc
				tupleOrderPrevious = gcpag.Asc
				cmpSignalTuple = CmpSignalLess
				cmpSignalPrevious = CmpSignalGreaterOrEqual
			}
			if pagReq.Order == gcpag.Desc {
				tupleOrder = tupleOrder.Reverse()
				tupleOrderPrevious = tupleOrderPrevious.Reverse()
				cmpSignalTuple = cmpSignalTuple.Reverse()
				cmpSignalPrevious = cmpSignalPrevious.Reverse()
			}
			return PagingCriteria[I]{
				Idt: idtValue,
				Query: fmt.Sprintf(
					`select * from (%s) as sub_gorepo_pag
					and (sub_gorepo_pag.%s, sub_gorepo_pag.%s) %s ((%s), %s)
					order by sub_gorepo_pag.%s %s, sub_gorepo_pag.%s %s
					%s
					union all
					select *, false as is_previous_page from (%s) as sub_gorepo_pag
					and (sub_gorepo_pag.%s, sub_gorepo_pag.%s) %s ((%s), %s)
					order by sub_gorepo_pag.%s %s, sub_gorepo_pag.%s %s
					%s`,
					query,
					sortFieldCriteria.ColumnName, idtColumn.Column, cmpSignalPrevious, sortFieldCriteria.SubQuery, syntax.PlaceHolder(idtColumn.PlaceHolder+1),
					sortFieldCriteria.ColumnName, tupleOrderPrevious.String(), idtColumn.Column, tupleOrderPrevious.String(),
					syntax.LimitStatement(idtColumn.PlaceHolder+2),
					query,
					sortFieldCriteria.ColumnName, idtColumn.Column, cmpSignalTuple, sortFieldCriteria.SubQuery, syntax.PlaceHolder(idtColumn.PlaceHolder+1),
					sortFieldCriteria.ColumnName, tupleOrder.String(), idtColumn.Column, tupleOrder.String(),
					syntax.LimitStatement(idtColumn.PlaceHolder+2)),
				Args: []any{idtValue, idtValue, 1, idtValue, idtValue, sizeMoreNextElement},
			}
		}
	}
	// With page parameter
	var (
		cmpSignal         CmpSign
		cmpSignalPrevious CmpSign
		idtOrder          gcpag.Order
	)
	switch pagReq.Orientation {
	case gcpag.NextPage:
		idtOrder = gcpag.Asc
		switch pagReq.StartPosition {
		case gcpag.StartAt:
			cmpSignal = CmpSignalGreaterOrEqual
			cmpSignalPrevious = CmpSignalLess
		case gcpag.AfterAt:
			cmpSignal = CmpSignalGreater
			cmpSignalPrevious = CmpSignalLessOrEqual
		}
	case gcpag.PreviousPage:
		idtOrder = gcpag.Desc
		switch pagReq.StartPosition {
		case gcpag.StartAt:
			cmpSignal = CmpSignalLessOrEqual
			cmpSignalPrevious = CmpSignalGreater
		case gcpag.AfterAt:
			cmpSignal = CmpSignalLess
			cmpSignalPrevious = CmpSignalGreaterOrEqual
		}
	}
	if pagReq.Order == gcpag.Desc {
		cmpSignal = cmpSignal.Reverse()
		cmpSignalPrevious = cmpSignalPrevious.Reverse()
		idtOrder = idtOrder.Reverse()
	}
	return PagingCriteria[I]{
		Idt: idtValue,
		Query: fmt.Sprintf(
			`(select *, true as is_previous_page from (%s) as sub_gorepo_pag
			sub_gorepo_pag.%s %s sub_gorepo_pag.%s
			order by sub_gorepo_pag.%s %s
			%s)
			union all
			(select *, false as is_previous_page from (%s) as sub_gorepo_pag
			sub_gorepo_pag.%s %s sub_gorepo_pag.%s
			order by sub_gorepo_pag.%s %s
			%s)`,
			query,
			idtColumn.Column, cmpSignalPrevious, syntax.PlaceHolder(idtColumn.PlaceHolder),
			idtColumn.Column, idtOrder.Reverse().String(),
			syntax.LimitStatement(idtColumn.PlaceHolder+1),
			query,
			idtColumn.Column, cmpSignal, syntax.PlaceHolder(idtColumn.PlaceHolder),
			idtColumn.Column, idtOrder.String(),
			syntax.LimitStatement(idtColumn.PlaceHolder+1)),

		Args: []any{idtValue, 1, idtValue, sizeMoreNextElement},
	}
}

func QueryPaginated[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](ctx context.Context, executor SqlExecutor, pagReq gcpag.PaginatedRequest[I], req QueryPaginatedParams[I, E]) (gcpag.PaginatedResponse[I, E], error) {
	var res gcpag.PaginatedResponse[I, E]
	syntax := executor.Syntax()
	pagCrit := newPagingCriteria(syntax, pagReq, req.QueryItems, ColumnCriteria{Column: req.IdtColumn, PlaceHolder: req.LastQueryPlaceHolder + 1}, req.FieldColumn)

	args := append(req.QueryArgs, pagCrit.Args...)
	fmt.Println("=======", pagCrit.Query) //TODO: REMOVER
	rowsItens, err := executor.Query(ctx, pagCrit.Query, args...)
	if err != nil {
		return res, fmt.Errorf("gcrepo: query itens paginated failed. Cause %w", err)
	}
	defer rowsItens.Close()
	var items []E
	hasPrevious := false
	if rowsItens.Next() {
		// check if previous row exists. If not
		rowsItemsWrapper := sqlRowWrapper{
			SqlRow: rowsItens,
			fnWrapper: func(args ...any) ([]any, error) {
				return append(args, &hasPrevious), nil
			},
		}
		m, err := req.ConverterQueryItems(&rowsItemsWrapper)
		if err != nil {
			return res, fmt.Errorf("gcrepo: convert item of the previous row failed. Cause %w", err)
		} else if !hasPrevious { // if not exist a previous, so it is page items
			items = append(items, m)
		}
		// continue taking page items
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
	}
	hasNextPage := len(items) > pagReq.Size
	if hasNextPage { //check if exist next row
		items = items[:len(items)-1]
	}
	if pagReq.Orientation == gcpag.PreviousPage {
		slices.Reverse(items)
	}
	return newPaginatedResponse(pagReq, items, hasPrevious, hasNextPage), nil
}

func newPaginatedResponse[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](pageReq gcpag.PaginatedRequest[I], items []E, hasPreviousPage, hasNextPage bool) gcpag.PaginatedResponse[I, E] {
	pageRes := gcpag.PaginatedResponse[I, E]{
		Items: items,
		Size:  pageReq.Size,
		Field: pageReq.SortField,
	}
	if len(items) > 0 {
		if hasPreviousPage {
			pageRes.PreviousPage = gcopt.Of(gcpag.AnotherPageRequest[I]{
				Idt:           items[0].Idt(),
				StartPosition: gcpag.AfterAt,
				Order:         pageReq.Order,
				Orientation:   gcpag.PreviousPage,
			})
		}
		if hasNextPage {
			pageRes.NextPage = gcopt.Of(gcpag.AnotherPageRequest[I]{
				Idt:           items[len(items)-1].Idt(),
				StartPosition: gcpag.AfterAt,
				Order:         pageReq.Order,
				Orientation:   gcpag.NextPage,
			})
		}
	}
	return pageRes
}
