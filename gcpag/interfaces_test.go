package gcpag

import "strconv"

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

type testUsersPagination struct {
	PaginatedRequest[testUserIdt]
}
