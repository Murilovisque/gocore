package gcpag

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/Murilovisque/gocore/gcopt"
)

type PaginatedRequest[I IdtOrdered] struct {
	Idt           gcopt.Optional[I]
	StartPosition StartPosition
	Order         Order
	Orientation   Orientation
	Size          int
	// Field gcopt.Optional[Field]
}

type AnotherPageRequest[I IdtOrdered] struct {
	Idt           I
	StartPosition StartPosition
	Order         Order
	Orientation   Orientation
	// Field gcopt.Optional[Field]
}

type PaginatedResponse[I IdtOrdered, M ModelOrderable[I]] struct {
	Items        []M
	SelfPage     gcopt.Optional[AnotherPageRequest[I]]
	NextPage     gcopt.Optional[AnotherPageRequest[I]]
	PreviousPage gcopt.Optional[AnotherPageRequest[I]]
	FirstPage    gcopt.Optional[AnotherPageRequest[I]]
	LastPage     gcopt.Optional[AnotherPageRequest[I]]
	Size         int
}

func (pageRes PaginatedResponse[I, M]) BuildHttpHeaderLinkValues(relativePath string) []string {
	links := make([]string, 0, 5)
	links = append(links, buildHttpHeaderRelatedLink(pageRes.FirstPage, "first", httpFieldPageStartIdt, httpFieldPageAfterIdt, relativePath, pageRes.Size)...)
	links = append(links, buildHttpHeaderRelatedLink(pageRes.SelfPage, "self", httpFieldPageStartIdt, httpFieldPageAfterIdt, relativePath, pageRes.Size)...)
	links = append(links, buildHttpHeaderRelatedLink(pageRes.NextPage, "next", httpFieldPageStartIdt, httpFieldPageAfterIdt, relativePath, pageRes.Size)...)
	links = append(links, buildHttpHeaderRelatedLink(pageRes.PreviousPage, "prev", httpFieldReversePageStartIdt, httpFieldReversePageAfterIdt, relativePath, pageRes.Size)...)
	links = append(links, buildHttpHeaderRelatedLink(pageRes.LastPage, "last", httpFieldReversePageStartIdt, httpFieldReversePageAfterIdt, relativePath, pageRes.Size)...)
	return links
}

func buildHttpHeaderRelatedLink[I IdtOrdered](page gcopt.Optional[AnotherPageRequest[I]], relation, fieldStartIdt, fieldAfterIdt, relativePath string, pageSize int) []string {
	anotherPage, ok := page.Take()
	if !ok {
		return []string{}
	}
	u, err := url.Parse(relativePath)
	if err != nil {
		slog.Default().Error("gcpag: invalid relative path", "path", relativePath)
		return []string{}
	}
	q := u.Query()
	switch anotherPage.StartPosition {
	case StartAt:
		q.Set(fieldStartIdt, anotherPage.Idt.String())
	case AfterAt:
		q.Set(fieldAfterIdt, anotherPage.Idt.String())
	default:
		slog.Default().Error("gcpag: invalid start position", "startPosition", anotherPage.StartPosition)
		return []string{}
	}
	q.Set(httpFieldPageSize, strconv.Itoa(pageSize))
	q.Set(httpFieldPageOrder, anotherPage.Order.String())
	u.RawQuery = q.Encode()
	return []string{fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relation)}
}
