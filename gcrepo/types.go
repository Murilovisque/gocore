package gcrepo

type CmpSign string

func (cp CmpSign) Reverse() CmpSign {
	switch cp {
	case CmpSignalGreater:
		return CmpSignalLess
	case CmpSignalLess:
		return CmpSignalGreater
	case CmpSignalGreaterOrEqual:
		return CmpSignalLessOrEqual
	case CmpSignalLessOrEqual:
		return CmpSignalGreaterOrEqual
	default:
		return ""
	}
}
