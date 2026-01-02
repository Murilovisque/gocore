package gcpag

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Murilovisque/gocore/gcopt"
)

func BuildHttpRequest[T IdtOrdered](req *http.Request, defaultOrder Order, defaultSize int, fnIdtParser func(string) (T, bool)) PaginatedRequest[T] {
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

func BuildResponse[T IdtOrdered, M ModelOrderable[T]](pageReq PaginatedRequest[T], itens []M, firstIdt, lastIdt T) PaginatedResponse[T, M] {
	if firstIdt > lastIdt {
		slog.Default().Warn("gcpag: last idt must be greather or equals to first idt", "first", firstIdt, "last", lastIdt)
		itens = []M{}
	}
	pageRes := PaginatedResponse[T, M]{
		Items: itens,
		Size:  pageReq.Size,
	}
	if len(itens) == 0 {
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
	initalItem := itens[0]
	if initalItem.OrderableIdt() > firstIdt {
		pageRes.PreviousPage = gcopt.Of(AnotherPageRequest[T]{
			Idt:           initalItem.OrderableIdt(),
			StartPosition: AfterAt,
			Order:         pageReq.Order,
			Orientation:   pageReq.Orientation.Reverse(),
		})
	}
	lastItem := itens[len(itens)-1]
	if lastItem.OrderableIdt() < lastIdt {
		pageRes.NextPage = gcopt.Of(AnotherPageRequest[T]{
			Idt:           lastItem.OrderableIdt(),
			StartPosition: AfterAt,
			Order:         pageReq.Order,
			Orientation:   pageReq.Orientation,
		})
	}
	pageRes.SelfPage = gcopt.Of(AnotherPageRequest[T]{
		Idt:           itens[0].OrderableIdt(),
		StartPosition: StartAt,
		Order:         pageReq.Order,
		Orientation:   pageReq.Orientation,
	})

	return pageRes
}
