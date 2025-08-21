package historic

type InformationType int

const (
	UnknownInformationType InformationType = iota
	Open
	High
	Low
	Close
	AdjustedClose
	Volume
	DividendAmount
	SplitCoefficient
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
	case Volume:
		return "volumne"
	case AdjustedClose:
		return "adjusted close"
	case DividendAmount:
		return "dividend amount"
	case SplitCoefficient:
		return "split coefficient"
	default:
		panic("invalid value of InformationType")
	}
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
	case Volume:
		return "5. volume"
	case AdjustedClose:
		return "6. adjusted close"
	case DividendAmount:
		return "7. dividend amount"
	case SplitCoefficient:
		return "8. split coefficient"
	default:
		panic("invalid value of InformationType")
	}
}
