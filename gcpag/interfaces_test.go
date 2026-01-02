package gcpag

import "strconv"

type testUserIdt int

func (t testUserIdt) String() string {
	return strconv.Itoa(int(t))
}

type testUserModel struct {
	Idt testUserIdt
}

func (t testUserModel) OrderableIdt() testUserIdt {
	return t.Idt
}

type testUsersPagination struct {
	PaginatedRequest[testUserIdt]
}
