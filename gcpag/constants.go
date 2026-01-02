package gcpag

const (
	AfterAt StartPosition = iota
	StartAt StartPosition = iota
)

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

const (
	NextPage     Orientation = iota
	PreviousPage Orientation = iota
)

const (
	httpFieldPageSize  = "page-size"
	httpFieldPageOrder = "page-order"
	// httpFieldPageBeforeIdt        = "page-before-idt"
	httpFieldPageStartIdt        = "page-start-idt"
	httpFieldPageAfterIdt        = "page-after-idt"
	httpFieldReversePageStartIdt = "reverse-page-start-idt"
	httpFieldReversePageAfterIdt = "reverse-page-after-idt"
	// httpFieldReversePageBeforeIdt = "reverse-page-before-idt"
)
