package gcpag

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"testing"

	"github.com/Murilovisque/gocore/gcopt"
)

func TestPaginationResponseBuildHttpHeaderLink(t *testing.T) {
	// no links
	pg := PaginatedResponse[testUserIdt, testUserModel]{
		Items: []testUserModel{
			{Idt: 1},
		},
		Size: 2,
	}
	result := pg.BuildHttpHeaderLinkValues("/users")
	if len(result) > 0 {
		t.Fatalf("expected empty but '%s'", result)
	}

	pg = PaginatedResponse[testUserIdt, testUserModel]{
		Items: []testUserModel{
			{Idt: 2},
		},
		Size: 1,
		FirstPage: gcopt.Of(AnotherPageRequest[testUserIdt]{
			Idt:           1,
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
		}),
		SelfPage: gcopt.Of(AnotherPageRequest[testUserIdt]{
			Idt:           2,
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
		}),
		NextPage: gcopt.Of(AnotherPageRequest[testUserIdt]{
			Idt:           2,
			StartPosition: AfterAt,
			Order:         Asc,
			Orientation:   NextPage,
		}),
		PreviousPage: gcopt.Of(AnotherPageRequest[testUserIdt]{
			Idt:           2,
			StartPosition: AfterAt,
			Order:         Asc,
			Orientation:   PreviousPage,
		}),
		LastPage: gcopt.Of(AnotherPageRequest[testUserIdt]{
			Idt:           3,
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   PreviousPage,
		}),
	}
	result = pg.BuildHttpHeaderLinkValues("/users")
	generateUrl := func(path, relation, field, idt, order string, size int) string {
		u, err := url.Parse(path)
		if err != nil {
			t.Fatal(err)
		}
		q := u.Query()
		q.Set(field, idt)
		q.Set(httpFieldPageSize, strconv.Itoa(size))
		q.Set(httpFieldPageOrder, order)
		u.RawQuery = q.Encode()
		return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relation)
	}
	expected := []string{
		generateUrl("/users", "first", httpFieldPageStartIdt, (pg.Items[0].Idt - 1).String(), pg.FirstPage.MustTake().Order.String(), pg.Size),
		generateUrl("/users", "self", httpFieldPageStartIdt, (pg.Items[0].Idt).String(), pg.SelfPage.MustTake().Order.String(), pg.Size),
		generateUrl("/users", "next", httpFieldPageAfterIdt, (pg.Items[0].Idt).String(), pg.NextPage.MustTake().Order.String(), pg.Size),
		generateUrl("/users", "prev", httpFieldReversePageAfterIdt, (pg.Items[0].Idt).String(), pg.PreviousPage.MustTake().Order.String(), pg.Size),
		generateUrl("/users", "last", httpFieldReversePageStartIdt, (pg.Items[0].Idt + 1).String(), pg.LastPage.MustTake().Order.String(), pg.Size),
	}
	if !slices.Equal(result, expected) {
		t.Fatalf("expected '%v' but '%v'", expected, result)
	}
}

// func TestPaginationRequest(t *testing.T) {
// 	var p testUsersPagination
// 	if p.Idt.IsPresent() {
// 		t.Fatalf("expected not present")
// 	}
// 	if p.Size != defaultSize {
// 		t.Fatalf("expected '%d', but '%d'", defaultSize, p.Size)
// 	}
// 	if p.Order != Asc {
// 		t.Fatalf("expected '%d', but '%d'", Asc, p.Order)
// 	}
// 	if p.Orientation != NextPage {
// 		t.Fatalf("expected '%d', but '%d'", NextPage, p.Orientation)
// 	}
// 	if p.StartPosition != AfterAt {
// 		t.Fatalf("expected '%d', but '%d'", AfterAt, p.StartPosition)
// 	}
// }
