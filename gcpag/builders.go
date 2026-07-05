package gcpag

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

func ParseRequestFromHttp[I gcfield.IdtOrdered](req *http.Request, params ParseRequestFromHttpParams[I]) (PaginatedRequest[I], error) {
	q := req.URL.Query()
	size, err := strconv.Atoi(q.Get(httpParamPageSize))
	if size < 1 || err != nil {
		size = params.DefaultSize
	}
	order, err := ParseOrder(q.Get(httpParamPageOrder))
	if err != nil {
		order = params.DefaultOrder
	}
	p := PaginatedRequest[I]{Size: size, Order: order}
	checkPoints := []struct {
		field string
		pos   StartPosition
		ori   Orientation
	}{
		{httpParamPageStartIdt, StartAt, NextPage},
		{httpParamPageAfterIdt, AfterAt, NextPage},
		{httpParamReversePageStartIdt, StartAt, PreviousPage},
		{httpParamReversePageAfterIdt, AfterAt, PreviousPage},
	}
	for _, cp := range checkPoints {
		if val := q.Get(cp.field); val != "" {
			if idt, ok, err := params.IdtParser(val); err != nil {
				return p, err
			} else if ok {
				p.Idt = gcopt.Of(idt)
				p.StartPosition = cp.pos
				p.Orientation = cp.ori
				break
			}
		}
	}
	if fp, ok := params.Field.Take(); ok {
		for _, name := range fp.AllowedNames() {
			if val := q.Get(httpParamPageSortField + name); val != "" {
				if ok, err := fp.Parse(name, val); err != nil {
					return p, err
				} else if ok {
					break
				}
			}
		}
	}
	return p, nil
}

// func NewResponse[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](pageReq PaginatedRequest[I], items []E, firstIdt, lastIdt gcopt.Optional[I]) PaginatedResponse[I, E] {
// 	pageRes := PaginatedResponse[I, E]{
// 		Items: items,
// 		Size:  pageReq.Size,
// 		Field: pageReq.Field,
// 	}
// 	if pageReq.Order == Desc {
// 		firstIdt, lastIdt = lastIdt, firstIdt
// 	}
// 	hasItems := len(items) > 0
// 	if fi, ok := firstIdt.Take(); ok {
// 		pageRes.FirstPage = gcopt.Of(AnotherPageRequest[I]{
// 			Idt:           fi,
// 			StartPosition: StartAt,
// 			Order:         pageReq.Order,
// 			Orientation:   NextPage,
// 		})
// 		if hasItems {
// 			initialItem := items[0]
// 			if (pageReq.Order == Asc && initialItem.Idt() > fi) || (pageReq.Order == Desc && initialItem.Idt() < fi) {
// 				pageRes.PreviousPage = gcopt.Of(AnotherPageRequest[I]{
// 					Idt:           initialItem.Idt(),
// 					StartPosition: AfterAt,
// 					Order:         pageReq.Order,
// 					Orientation:   PreviousPage,
// 				})
// 			}
// 		}
// 	}
// 	if li, ok := lastIdt.Take(); ok {
// 		pageRes.LastPage = gcopt.Of(AnotherPageRequest[I]{
// 			Idt:           li,
// 			StartPosition: StartAt,
// 			Order:         pageReq.Order,
// 			Orientation:   PreviousPage,
// 		})
// 		if hasItems {
// 			lastItem := items[len(items)-1]
// 			if (pageReq.Order == Asc && lastItem.Idt() < li) || (pageReq.Order == Desc && lastItem.Idt() > li) {
// 				pageRes.NextPage = gcopt.Of(AnotherPageRequest[I]{
// 					Idt:           lastItem.Idt(),
// 					StartPosition: AfterAt,
// 					Order:         pageReq.Order,
// 					Orientation:   NextPage,
// 				})
// 			}
// 		}
// 	}
// 	if hasItems {
// 		pageRes.SelfPage = gcopt.Of(AnotherPageRequest[I]{
// 			Idt:           items[0].Idt(),
// 			StartPosition: StartAt,
// 			Order:         pageReq.Order,
// 			Orientation:   pageReq.Orientation,
// 		})
// 	}
// 	return pageRes
// }

func ParseOrder(vl string) (Order, error) {
	vl = strings.ToLower(vl)
	switch vl {
	case "asc":
		return Asc, nil
	case "desc":
		return Desc, nil
	default:
		return Asc, fmt.Errorf("gcpag: invalid order value '%s'", vl)
	}
}
