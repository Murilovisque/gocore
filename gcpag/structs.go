package gcpag

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

type PaginatedRequest[I gcfield.IdtOrdered] struct {
	Idt           gcopt.Optional[I]
	StartPosition StartPosition
	Order         Order
	Orientation   Orientation
	Size          int
	Field         gcopt.Optional[gcfield.FieldParser]
}

type AnotherPageRequest[I gcfield.IdtOrdered] struct {
	Idt           I
	StartPosition StartPosition
	Order         Order
	Orientation   Orientation
}

type PaginatedResponse[I gcfield.IdtOrdered, E gcfield.Identifiable[I]] struct {
	Items        []E
	SelfPage     gcopt.Optional[AnotherPageRequest[I]]
	NextPage     gcopt.Optional[AnotherPageRequest[I]]
	PreviousPage gcopt.Optional[AnotherPageRequest[I]]
	FirstPage    gcopt.Optional[AnotherPageRequest[I]]
	LastPage     gcopt.Optional[AnotherPageRequest[I]]
	Size         int
	Field        gcopt.Optional[gcfield.FieldParser]
}

func (pageRes PaginatedResponse[I, E]) ParseHttpHeaderLinkValues(relativePath string) ([]string, error) {
	links := make([]string, 0, 5)
	addLink := func(page gcopt.Optional[AnotherPageRequest[I]], relation, fieldStartIdt, fieldAfterIdt string) error {
		anotherPage, ok := page.Take()
		if !ok {
			return nil
		}
		u, err := url.Parse(relativePath)
		if err != nil {
			return fmt.Errorf("gcpag: invalid relative path '%s', http header related link build failed. Cause %w", relativePath, err)
		}
		q := u.Query()
		switch anotherPage.StartPosition {
		case StartAt:
			q.Set(fieldStartIdt, anotherPage.Idt.String())
		case AfterAt:
			q.Set(fieldAfterIdt, anotherPage.Idt.String())
		default:
			return fmt.Errorf("gcpag: invalid start position '%s', http header related link build failed. Cause %w", anotherPage.StartPosition, err)
		}
		q.Set(httpParamPageSize, strconv.Itoa(pageRes.Size))
		q.Set(httpParamPageOrder, anotherPage.Order.String())
		if f, ok := pageRes.Field.Take(); ok && f.IsValid() {
			q.Set(httpParamPageField+f.Name(), f.String())
		}
		u.RawQuery = q.Encode()
		links = append(links, []string{fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relation)}...)
		return nil
	}
	var err error
	if err = addLink(pageRes.FirstPage, "first", httpParamPageStartIdt, httpParamPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.SelfPage, "self", httpParamPageStartIdt, httpParamPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.NextPage, "next", httpParamPageStartIdt, httpParamPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.PreviousPage, "prev", httpParamReversePageStartIdt, httpParamReversePageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.LastPage, "last", httpParamReversePageStartIdt, httpParamReversePageAfterIdt); err != nil {
		return nil, err
	}
	return links, nil
}

func (pageRes PaginatedResponse[I, E]) SelfPageAsRequest() gcopt.Optional[PaginatedRequest[I]] {
	return pageRes.pageAsRequest(pageRes.SelfPage)
}

func (pageRes PaginatedResponse[I, E]) FirstPageAsRequest() gcopt.Optional[PaginatedRequest[I]] {
	return pageRes.pageAsRequest(pageRes.FirstPage)
}

func (pageRes PaginatedResponse[I, E]) LastPageAsRequest() gcopt.Optional[PaginatedRequest[I]] {
	return pageRes.pageAsRequest(pageRes.LastPage)
}

func (pageRes PaginatedResponse[I, E]) NextPageAsRequest() gcopt.Optional[PaginatedRequest[I]] {
	return pageRes.pageAsRequest(pageRes.NextPage)
}

func (pageRes PaginatedResponse[I, E]) PreviousPageAsRequest() gcopt.Optional[PaginatedRequest[I]] {
	return pageRes.pageAsRequest(pageRes.PreviousPage)
}

func (pageRes PaginatedResponse[I, E]) pageAsRequest(page gcopt.Optional[AnotherPageRequest[I]]) gcopt.Optional[PaginatedRequest[I]] {
	p, ok := page.Take()
	if !ok {
		return gcopt.Empty[PaginatedRequest[I]]()
	}
	return gcopt.Of(PaginatedRequest[I]{
		Idt:           gcopt.Of(p.Idt),
		StartPosition: p.StartPosition,
		Order:         p.Order,
		Orientation:   p.Orientation,
		Size:          pageRes.Size,
		Field:         pageRes.Field,
	})
}
