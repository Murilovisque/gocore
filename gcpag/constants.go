package gcpag

const (
	AfterAt StartPosition = iota
	StartAt
)

const (
	Asc Order = iota
	Desc
)

const (
	NextPage Orientation = iota
	PreviousPage
)

const (
	httpParamPageSize            = "page-size"
	httpParamPageSortField       = "page-sort-field"
	httpParamPageOrder           = "page-order"
	httpParamPageStartIdt        = "page-start-idt"
	httpParamPageAfterIdt        = "page-after-idt"
	httpParamReversePageStartIdt = "reverse-page-start-idt"
	httpParamReversePageAfterIdt = "reverse-page-after-idt"
)

// httpParamPageBeforeIdt        = "page-before-idt"
// httpFieldReversePageBeforeIdt = "reverse-page-before-idt"
