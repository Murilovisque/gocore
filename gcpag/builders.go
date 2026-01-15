package gcpag

import (
	"net/http"
	"strconv"

	"github.com/Murilovisque/gocore/gcid"
	"github.com/Murilovisque/gocore/gcopt"
)

func BuildByHttpRequest[T gcid.IdtOrdered](req *http.Request, defaultOrder Order, defaultSize int, fnIdtParser func(string) (T, bool)) PaginatedRequest[T] {
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
		// } else if idt, valid = idtConverter(q.Get(httpFieldPageBeforeIdt)); valid {
		// 	p.Idt = gcopt.Of(idt)
		// 	p.StartPosition = AfterAt
		// 	p.Orientation = PreviousPage
		// 	p.Order = Asc
	} else if idt, valid = fnIdtParser(q.Get(httpFieldReversePageStartIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = StartAt
		p.Orientation = NextPage
	} else if idt, valid = fnIdtParser(q.Get(httpFieldReversePageAfterIdt)); valid {
		p.Idt = gcopt.Of(idt)
		p.StartPosition = AfterAt
		p.Orientation = NextPage
		// } else if idt, valid = idtConverter(q.Get(httpFieldReversePageBeforeIdt)); valid {
		// 	p.Idt = gcopt.Of(idt)
		// 	p.StartPosition = AfterAt
		// 	p.Orientation = PreviousPage
		// 	p.Order = Desc
	}
	return p
}

func BuildResponse[T gcid.IdtOrdered, M gcid.Identifiable[T]](pageReq PaginatedRequest[T], items []M, firstIdt, lastIdt T) PaginatedResponse[T, M] {
	// if firstIdt > lastIdt {
	// 	slog.Default().Warn("gcpag: last idt must be greather or equals to first idt", "first", firstIdt, "last", lastIdt)
	// 	itens = []M{}
	// }
	pageRes := PaginatedResponse[T, M]{
		Items: items,
		Size:  pageReq.Size,
	}
	if len(items) == 0 {
		return pageRes
	}
	pageRes.FirstPage = gcopt.Of(AnotherPageRequest[T]{
		Idt:           firstIdt,
		StartPosition: StartAt,
		Order:         pageReq.Order,
		Orientation:   pageReq.Orientation,
	})
	pageRes.LastPage = gcopt.Of(AnotherPageRequest[T]{
		Idt:           lastIdt,
		StartPosition: StartAt,
		Order:         pageReq.Order,
		Orientation:   pageReq.Orientation.Reverse(),
	})
	initalItem := items[0]
	if initalItem.Idt() > firstIdt {
		pageRes.PreviousPage = gcopt.Of(AnotherPageRequest[T]{
			Idt:           initalItem.Idt(),
			StartPosition: AfterAt,
			Order:         pageReq.Order,
			Orientation:   pageReq.Orientation.Reverse(),
		})
	}
	lastItem := items[len(items)-1]
	if lastItem.Idt() < lastIdt {
		pageRes.NextPage = gcopt.Of(AnotherPageRequest[T]{
			Idt:           lastItem.Idt(),
			StartPosition: AfterAt,
			Order:         pageReq.Order,
			Orientation:   pageReq.Orientation,
		})
	}
	pageRes.SelfPage = gcopt.Of(AnotherPageRequest[T]{
		Idt:           items[0].Idt(),
		StartPosition: StartAt,
		Order:         pageReq.Order,
		Orientation:   pageReq.Orientation,
	})

	return pageRes
}
