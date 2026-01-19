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
	// Field gcopt.Optional[Field]
}

type AnotherPageRequest[I gcfield.IdtOrdered] struct {
	Idt           I
	StartPosition StartPosition
	Order         Order
	Orientation   Orientation
	// Field gcopt.Optional[Field]
}

func (a AnotherPageRequest[I]) ToPageRequest(size int) PaginatedRequest[I] {
	return PaginatedRequest[I]{
		Idt:           gcopt.Of(a.Idt),
		StartPosition: a.StartPosition,
		Order:         a.Order,
		Orientation:   a.Orientation,
		Size:          size,
	}
}

type PaginatedResponse[I gcfield.IdtOrdered, E gcfield.Identifiable[I]] struct {
	Items        []E
	SelfPage     gcopt.Optional[AnotherPageRequest[I]]
	NextPage     gcopt.Optional[AnotherPageRequest[I]]
	PreviousPage gcopt.Optional[AnotherPageRequest[I]]
	FirstPage    gcopt.Optional[AnotherPageRequest[I]]
	LastPage     gcopt.Optional[AnotherPageRequest[I]]
	Size         int
}

func (pageRes PaginatedResponse[I, E]) BuildHttpHeaderLinkValues(relativePath string) ([]string, error) {
	links := make([]string, 0, 5)
	addLink := func(page gcopt.Optional[AnotherPageRequest[I]], relation, fieldStartIdt, fieldAfterIdt string) error {
		anotherPage, ok := page.Take()
		if !ok {
			return nil
		}
		u, err := url.Parse(relativePath)
		if err != nil {
			return fmt.Errorf("gcpag: invalid relative path '%s', http header releted link build failed. Cause %w", relativePath, err)
		}
		q := u.Query()
		switch anotherPage.StartPosition {
		case StartAt:
			q.Set(fieldStartIdt, anotherPage.Idt.String())
		case AfterAt:
			q.Set(fieldAfterIdt, anotherPage.Idt.String())
		default:
			return fmt.Errorf("gcpag: invalid start position '%s', http header releted link build failed. Cause %w", anotherPage.StartPosition, err)
		}
		q.Set(httpFieldPageSize, strconv.Itoa(pageRes.Size))
		q.Set(httpFieldPageOrder, anotherPage.Order.String())
		u.RawQuery = q.Encode()
		links = append(links, []string{fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relation)}...)
		return nil
	}
	var err error
	if err = addLink(pageRes.FirstPage, "first", httpFieldPageStartIdt, httpFieldPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.SelfPage, "self", httpFieldPageStartIdt, httpFieldPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.NextPage, "next", httpFieldPageStartIdt, httpFieldPageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.PreviousPage, "prev", httpFieldReversePageStartIdt, httpFieldReversePageAfterIdt); err != nil {
		return nil, err
	}
	if err = addLink(pageRes.LastPage, "last", httpFieldReversePageStartIdt, httpFieldReversePageAfterIdt); err != nil {
		return nil, err
	}
	return links, nil
}
