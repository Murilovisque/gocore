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

	setupFirstPage := func(items []testUserModel, pageSize int) {
		currentFirstItemPage = 0
		currentLastItemPage = pageSize - 1
		if currentLastItemPage >= len(items) {
			currentLastItemPage = len(items) - 1
		}
		t.Logf("current first: %d, current last: %d", currentFirstItemPage, currentLastItemPage)
	}
	setupNextPage := func(items []testUserModel, pageSize int) {
		currentFirstItemPage = currentLastItemPage + 1
		currentLastItemPage = currentFirstItemPage + pageSize - 1
		if currentLastItemPage >= len(items) {
			currentLastItemPage = len(items) - 1
		}
		t.Logf("current first: %d, current last: %d", currentFirstItemPage, currentLastItemPage)
	}
	doBuildResponse := func(items []testUserModel, pageSize, expectedItemsSliceLength int) PaginatedResponse[testUserIdt, testUserModel] {
		pgReq := PaginatedRequest[testUserIdt]{
			Idt:           gcopt.Of(items[currentFirstItemPage].Idt()),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   NextPage,
			Size:          pageSize,
		}
		itemsQuery := items[currentFirstItemPage : currentLastItemPage+1]
		res := BuildResponse(pgReq, itemsQuery, gcopt.Of(testUserIdt(0)), gcopt.Of(items[len(items)-1].Idt()))
		t.Logf("response built: %v", res)
		if len(res.Items) != expectedItemsSliceLength {
			t.Fatalf("expected length %d, but %d", expectedItemsSliceLength, len(res.Items))
		} else if !slices.Equal(itemsQuery, res.Items) {
			t.Fatalf("expected %v, but %v", itemsQuery, res.Items)
		} else if res.Size != pageSize {
			t.Fatalf("expected %d, but %d", pageSize, res.Size)
		}
		return res
	}
	assertSelfFirstAndLastPages := func(items []testUserModel, res PaginatedResponse[testUserIdt, testUserModel]) {
		// valid self page
		expectPage := AnotherPageRequest[testUserIdt]{
			Idt:           items[currentFirstItemPage].Idt(),
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
			Idt:           items[0].Idt(),
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
			Idt:           items[len(items)-1].Idt(),
			StartPosition: StartAt,
			Order:         Asc,
			Orientation:   PreviousPage,
		}
		pageResult = res.LastPage.MustTake()
		if !reflect.DeepEqual(pageResult, expectPage) {
			t.Fatalf("expected %v, but %v", expectPage, pageResult)
		}
	}
	assertNextPage := func(items []testUserModel, res PaginatedResponse[testUserIdt, testUserModel], exists bool) {
		// valid next
		if exists {
			expectPage := AnotherPageRequest[testUserIdt]{
				Idt:           items[currentLastItemPage].Idt(),
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
	assertPreviousPage := func(items []testUserModel, res PaginatedResponse[testUserIdt, testUserModel], exists bool) {
		// valid previous
		if exists {
			expectPage := AnotherPageRequest[testUserIdt]{
				Idt:           items[currentFirstItemPage].Idt(),
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
	items := []testUserModel{
		{
			idt: 0,
		},
		{
			idt: 1,
		},
		{
			idt: 2,
		},
		{
			idt: 3,
		},
		{
			idt: 4,
		},
	}
	// Test case
	{
		pageSize := 1
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, false)

		// Test next pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, true)

		// Test next pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, true)

		// Test next pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, true)

		// Test last pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, true)
	}
	// Test case
	{
		pageSize := 2
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 2)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, false)

		// Test next pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 2)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, true)

		// Test last pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, true)
	}
	// Test case
	{
		pageSize := 3
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 3)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, false)

		// Test last pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 2)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, true)
	}

	// Test case
	{
		pageSize := 4
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 4)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, true)
		assertPreviousPage(items, res, false)

		// Test last pagination
		setupNextPage(items, pageSize)
		res = doBuildResponse(items, pageSize, 1)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, true)
	}

	// Test case
	{
		pageSize := 5
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first and last pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 5)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, false)
	}
	// Test case
	{
		pageSize := 6
		t.Logf("test case: {size: %d, itens: %v}", pageSize, items)
		// Test first and last pagination
		setupFirstPage(items, pageSize)
		res := doBuildResponse(items, pageSize, 5)
		assertSelfFirstAndLastPages(items, res)
		assertNextPage(items, res, false)
		assertPreviousPage(items, res, false)
	}
}
