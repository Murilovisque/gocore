package gcpag

import (
	"log/slog"
)

type StartPosition int

func (sp StartPosition) IsValid() bool {
	return sp == AfterAt || sp == StartAt
}

func (sp StartPosition) String() string {
	switch sp {
	case StartAt:
		return "startAt"
	case AfterAt:
		return "afterAt"
	default:
		return ""
	}
}

type Order string

func (o Order) Reverse() Order {
	switch o {
	case Asc:
		return Desc
	case Desc:
		return Asc
	default:
		slog.Default().Warn("gcpag: invalid page order, reverse failed", "invalidOrder", o)
		return o
	}
}

func (o Order) IsValid() bool {
	return o == Asc || o == Desc
}

func (o Order) String() string {
	return string(o)
}

type Orientation int

func (o Orientation) Reverse() Orientation {
	switch o {
	case NextPage:
		return PreviousPage
	case PreviousPage:
		return NextPage
	default:
		slog.Default().Warn("gcpag: invalid page orientation, reverse failed", "invalidOrder", o)
		return o
	}
}

func (o Orientation) IsValid() bool {
	return o == NextPage || o == PreviousPage
}

func (o Orientation) String() string {
	switch o {
	case NextPage:
		return "nextPage"
	case PreviousPage:
		return "previousPage"
	default:
		return ""
	}
}
