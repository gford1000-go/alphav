package fx

type InformationType int

const (
	UnknownInformationType InformationType = iota
	Open
	High
	Low
	Close
	InvalidInformationType
)

func (i InformationType) String() string {
	switch i {
	case Open:
		return "open"
	case High:
		return "high"
	case Low:
		return "low"
	case Close:
		return "close"
	default:
		panic("invalid value of InformationType")
	}
}

func (i InformationType) isValid() bool {
	if i <= UnknownInformationType || i >= InvalidInformationType {
		return false
	}
	return true
}

func (i InformationType) toAVString() string {
	switch i {
	case Open:
		return "1. open"
	case High:
		return "2. high"
	case Low:
		return "3. low"
	case Close:
		return "4. close"
	default:
		panic("invalid value of InformationType")
	}
}
