package gcpag

import (
	"net/http"
	"strconv"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

func BuildRequestFromHttp[T gcfield.IdtOrdered](req *http.Request, defaultOrder Order, defaultSize int, fnIdtParser func(string) (T, bool)) PaginatedRequest[T] {
	q := req.URL.Query()
	size, err := strconv.Atoi(q.Get(httpFieldPageSize))
	if size < 1 || err != nil {
		size = defaultSize
	}
	order, err := StringToOrder(q.Get(httpFieldPageOrder))
	if err != nil {
		order = defaultOrder
	}
	p := PaginatedRequest[T]{Size: size, Order: order}
	var idt T
	var valid bool
	if idt, valid = fnIdtParser(q.Get(httpFieldPageStartIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = StartAt
		p.Orientation = NextPage
	} else if idt, valid = fnIdtParser(q.Get(httpFieldPageAfterIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = AfterAt
		p.Orientation = NextPage
	} else if idt, valid = fnIdtParser(q.Get(httpFieldReversePageStartIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = StartAt
		p.Orientation = NextPage
	} else if idt, valid = fnIdtParser(q.Get(httpFieldReversePageAfterIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = AfterAt
		p.Orientation = NextPage
	}
	return p
}

func BuildResponse[I gcfield.IdtOrdered, E gcfield.Identifiable[I]](pageReq PaginatedRequest[I], items []E, firstIdt, lastIdt gcopt.Optional[I]) PaginatedResponse[I, E] {
	pageRes := PaginatedResponse[I, E]{
		Items: items,
		Size:  pageReq.Size,
	}
	if pageReq.Order == Desc {
		firstIdt, lastIdt = lastIdt, firstIdt
	}
	hasItems := len(items) > 0
	if fi, ok := firstIdt.Take(); ok {
		pageRes.FirstPage = gcopt.Of(AnotherPageRequest[I]{
			Idt:           fi,
			StartPosition: StartAt,
			Order:         pageReq.Order,
			Orientation:   NextPage,
		})
		if hasItems {
			initialItem := items[0]
			if (pageReq.Order == Asc && initialItem.Idt() > fi) || (pageReq.Order == Desc && initialItem.Idt() < fi) {
				pageRes.PreviousPage = gcopt.Of(AnotherPageRequest[I]{
					Idt:           initialItem.Idt(),
					StartPosition: AfterAt,
					Order:         pageReq.Order,
					Orientation:   PreviousPage,
				})
			}
		}
	}
	if li, ok := lastIdt.Take(); ok {
		pageRes.LastPage = gcopt.Of(AnotherPageRequest[I]{
			Idt:           li,
			StartPosition: StartAt,
			Order:         pageReq.Order,
			Orientation:   PreviousPage,
		})
		if hasItems {
			lastItem := items[len(items)-1]
			if (pageReq.Order == Asc && lastItem.Idt() < li) || (pageReq.Order == Desc && lastItem.Idt() > li) {
				pageRes.NextPage = gcopt.Of(AnotherPageRequest[I]{
					Idt:           lastItem.Idt(),
					StartPosition: AfterAt,
					Order:         pageReq.Order,
					Orientation:   NextPage,
				})
			}
		}
	}
	if hasItems {
		pageRes.SelfPage = gcopt.Of(AnotherPageRequest[I]{
			Idt:           items[0].Idt(),
			StartPosition: StartAt,
			Order:         pageReq.Order,
			Orientation:   pageReq.Orientation,
		})
	}
	return pageRes
}
