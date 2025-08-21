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
		return "volume"
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
	case AdjustedClose:
		return "5. adjusted close"
	case Volume:
		return "6. volume"
	case DividendAmount:
		return "7. dividend amount"
	case SplitCoefficient:
		return "8. split coefficient"
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
