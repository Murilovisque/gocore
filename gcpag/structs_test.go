package gcpag

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"testing"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

func TestPaginationResponseBuildHttpHeaderLink(t *testing.T) {
	// no links
	pg := PaginatedResponse[testUserIdt, testUserModel]{
		Items: []testUserModel{
			{idt: 1},
		},
		Size: 2,
	}
	result, err := pg.ParseHttpHeaderLinkValues("/users")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) > 0 {
		t.Fatalf("expected empty but '%s'", result)
	}

	pg = PaginatedResponse[testUserIdt, testUserModel]{
		Items: []testUserModel{
			{idt: 2},
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
	generateUrl := func(path, relation, field, idt, order string, size int, extraField gcopt.Optional[gcfield.FieldParser]) string {
		u, err := url.Parse(path)
		if err != nil {
			t.Fatal(err)
		}
		q := u.Query()
		q.Set(field, idt)
		q.Set(httpParamPageSize, strconv.Itoa(size))
		q.Set(httpParamPageOrder, order)
		if f, ok := extraField.Take(); ok && f.IsValid() {
			q.Set(httpParamPageField+f.Name(), f.String())
		}
		u.RawQuery = q.Encode()
		return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relation)
	}
	expected := []string{
		generateUrl("/users", "first", httpParamPageStartIdt, (pg.Items[0].Idt() - 1).String(), pg.FirstPage.MustTake().Order.String(), pg.Size, gcopt.Empty[gcfield.FieldParser]()),
		generateUrl("/users", "self", httpParamPageStartIdt, (pg.Items[0].Idt()).String(), pg.SelfPage.MustTake().Order.String(), pg.Size, gcopt.Empty[gcfield.FieldParser]()),
		generateUrl("/users", "next", httpParamPageAfterIdt, (pg.Items[0].Idt()).String(), pg.NextPage.MustTake().Order.String(), pg.Size, gcopt.Empty[gcfield.FieldParser]()),
		generateUrl("/users", "prev", httpParamReversePageAfterIdt, (pg.Items[0].Idt()).String(), pg.PreviousPage.MustTake().Order.String(), pg.Size, gcopt.Empty[gcfield.FieldParser]()),
		generateUrl("/users", "last", httpParamReversePageStartIdt, (pg.Items[0].Idt() + 1).String(), pg.LastPage.MustTake().Order.String(), pg.Size, gcopt.Empty[gcfield.FieldParser]()),
	}

	result, err = pg.ParseHttpHeaderLinkValues("/users")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(result, expected) {
		t.Fatalf("expected '%v' but '%v'", expected, result)
	}

	// extra Field paginated
	fieldParser := gcfield.NewFieldParser([]string{"name"}, func(name, value string) (parsedValue any, valid bool, err error) {
		valid = name == "name"
		if valid {
			parsedValue = value
		}
		return
	}, func(name string, value any) string {
		return value.(string)
	})
	if ok, err := fieldParser.Parse("name", "Murilo"); err != nil || !ok {
		t.Fatalf("set field failed: %t and %s", ok, err.Error())
	}
	pg.Field = gcopt.Of(fieldParser)
	expected = []string{
		generateUrl("/users", "first", httpParamPageStartIdt, (pg.Items[0].Idt() - 1).String(), pg.FirstPage.MustTake().Order.String(), pg.Size, gcopt.Of(fieldParser)),
		generateUrl("/users", "self", httpParamPageStartIdt, (pg.Items[0].Idt()).String(), pg.SelfPage.MustTake().Order.String(), pg.Size, gcopt.Of(fieldParser)),
		generateUrl("/users", "next", httpParamPageAfterIdt, (pg.Items[0].Idt()).String(), pg.NextPage.MustTake().Order.String(), pg.Size, gcopt.Of(fieldParser)),
		generateUrl("/users", "prev", httpParamReversePageAfterIdt, (pg.Items[0].Idt()).String(), pg.PreviousPage.MustTake().Order.String(), pg.Size, gcopt.Of(fieldParser)),
		generateUrl("/users", "last", httpParamReversePageStartIdt, (pg.Items[0].Idt() + 1).String(), pg.LastPage.MustTake().Order.String(), pg.Size, gcopt.Of(fieldParser)),
	}

	result, err = pg.ParseHttpHeaderLinkValues("/users")
	if err != nil {
		t.Fatal(err)
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

type testUserIdt int

func (t testUserIdt) String() string {
	return strconv.Itoa(int(t))
}

type testUserModel struct {
	idt testUserIdt
}

func (t testUserModel) Idt() testUserIdt {
	return t.idt
}

// type testUsersPagination struct {
// 	PaginatedRequest[testUserIdt]
// }
