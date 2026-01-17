package gcpag

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

type Order int

func (o Order) Reverse() Order {
	switch o {
	case Asc:
		return Desc
	case Desc:
		return Asc
	default:
		return o
	}
}

func (o Order) IsValid() bool {
	return o == Asc || o == Desc
}

func (o Order) String() string {
	switch o {
	case Asc:
		return "asc"
	case Desc:
		return "desc"
	default:
		return ""
	}
}

type Orientation int

func (o Orientation) Reverse() Orientation {
	switch o {
	case NextPage:
		return PreviousPage
	case PreviousPage:
		return NextPage
	default:
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
