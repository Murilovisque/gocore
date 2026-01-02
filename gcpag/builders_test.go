package gcpag

import (
	"reflect"
	"slices"
	"testing"

	"github.com/Murilovisque/gocore/gcopt"
)

func TestBuildResponse(t *testing.T) {
	// utils and test functions
	var currentFirstItemPage, currentLastItemPage int

	setupFirstPage := func(itens []testUserModel, pageSize int) {
		currentFirstItemPage = 0
		currentLastItemPage = pageSize - 1
		if currentLastItemPage >= len(itens) {
			currentLastItemPage = len(itens) - 1
		}
		t.Logf("current first: %d, current last: %d", currentFirstItemPage, currentLastItemPage)
	}
	setupNextPage := func(itens []testUserModel, pageSize int) {
		currentFirstItemPage = currentLastItemPage + 1
		currentLastItemPage = currentFirstItemPage + pageSize - 1
		if currentLastItemPage >= len(itens) {
			currentLastItemPage = len(itens) - 1
		}
		t.Logf("current first: %d, current last: %d", currentFirstItemPage, currentLastItemPage)
	}
	doBuildResponse := func(itens []testUserModel, pageSize, expectedItensSliceLength int) PaginatedResponse[testUserIdt, testUserModel] {
		pgReq := PaginatedRequest[testUserIdt]{
			Idt:           gcopt.Of(itens[currentFirstItemPage].OrderableIdt()),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
			Size:          pageSize,
		}
		itensQuery := itens[currentFirstItemPage : currentLastItemPage+1]
		res := BuildResponse(pgReq, itensQuery, 0, itens[len(itens)-1].Idt)
		t.Logf("response built: %v", res)
		if len(res.Items) != expectedItensSliceLength {
			t.Fatalf("expected length %d, but %d", expectedItensSliceLength, len(res.Items))
		} else if !slices.Equal(itensQuery, res.Items) {
			t.Fatalf("expected %v, but %v", itensQuery, res.Items)
		} else if res.Size != pageSize {
			t.Fatalf("expected %d, but %d", pageSize, res.Size)
		}
		return res
	}
	assertSelfFirstAndLastPages := func(itens []testUserModel, res PaginatedResponse[testUserIdt, testUserModel]) {
		// valid self page
		expectPage := AnotherPageRequest[testUserIdt]{
			Idt:           itens[currentFirstItemPage].OrderableIdt(),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
		}
		pageResult := res.SelfPage.MustTake()
		if !reflect.DeepEqual(pageResult, expectPage) {
			t.Fatalf("expected %v, but %v", expectPage, pageResult)
		}
		// valid first page
		expectPage = AnotherPageRequest[testUserIdt]{
			Idt:           itens[0].OrderableIdt(),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
		}
		pageResult = res.FirstPage.MustTake()
		if !reflect.DeepEqual(pageResult, expectPage) {
			t.Fatalf("expected %v, but %v", expectPage, pageResult)
		}
		// valid last page
		expectPage = AnotherPageRequest[testUserIdt]{
			Idt:           itens[len(itens)-1].OrderableIdt(),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   PreviousPage,
		}
		pageResult = res.LastPage.MustTake()
		if !reflect.DeepEqual(pageResult, expectPage) {
			t.Fatalf("expected %v, but %v", expectPage, pageResult)
		}
	}
	assertNextPage := func(itens []testUserModel, res PaginatedResponse[testUserIdt, testUserModel], exists bool) {
		// valid next
		if exists {
			expectPage := AnotherPageRequest[testUserIdt]{
				Idt:           itens[currentLastItemPage].OrderableIdt(),
				StartPosition: AfterAt,
				Order:         Asc,
				Orientation:   NextPage,
			}
			pageResult := res.NextPage.MustTake()
			if !reflect.DeepEqual(pageResult, expectPage) {
				t.Fatalf("expected %v, but %v", expectPage, pageResult)
			}
		} else if res.NextPage.IsPresent() {
			t.Fatalf("expected does not exists next page")
		}
	}
	assertPreviousPage := func(itens []testUserModel, res PaginatedResponse[testUserIdt, testUserModel], exists bool) {
		// valid previous
		if exists {
			expectPage := AnotherPageRequest[testUserIdt]{
				Idt:           itens[currentFirstItemPage].OrderableIdt(),
				StartPosition: AfterAt,
				Order:         Asc,
				Orientation:   PreviousPage,
			}
			pageResult := res.PreviousPage.MustTake()
			if !reflect.DeepEqual(pageResult, expectPage) {
				t.Fatalf("expected %v, but %v", expectPage, pageResult)
			}
		} else if res.PreviousPage.IsPresent() {
			t.Fatalf("expected does not exists previous page")
		}
	}
	// Test cases variables
	itens := []testUserModel{
		{
			Idt: 0,
		},
		{
			Idt: 1,
		},
		{
			Idt: 2,
		},
		{
			Idt: 3,
		},
		{
			Idt: 4,
		},
	}
	// Test case
	{
		pageSize := 1
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, false)

		// Test next pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, true)

		// Test next pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, true)

		// Test next pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, true)

		// Test last pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, true)
	}
	// Test case
	{
		pageSize := 2
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 2)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, false)

		// Test next pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 2)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, true)

		// Test last pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, true)
	}
	// Test case
	{
		pageSize := 3
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 3)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, false)

		// Test last pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 2)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, true)
	}

	// Test case
	{
		pageSize := 4
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 4)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, true)
		assertPreviousPage(itens, res, false)

		// Test last pagination
		setupNextPage(itens, pageSize)
		res = doBuildResponse(itens, pageSize, 1)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, true)
	}

	// Test case
	{
		pageSize := 5
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first and last pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 5)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, false)
	}
	// Test case
	{
		pageSize := 6
		t.Logf("test case: {size: %d, itens: %v}", pageSize, itens)
		// Test first and last pagination
		setupFirstPage(itens, pageSize)
		res := doBuildResponse(itens, pageSize, 5)
		assertSelfFirstAndLastPages(itens, res)
		assertNextPage(itens, res, false)
		assertPreviousPage(itens, res, false)
	}
}
