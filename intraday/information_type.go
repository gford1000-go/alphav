package intraday

type InformationType int

const (
	UnknownIntradayInformationType InformationType = iota
	Open
	High
	Low
	Close
	Volume
	InvalidIntradayInformationType
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
	default:
		panic("invalid value of IntradayInformationType")
	}
}

func (i InformationType) isValid() bool {
	if i <= UnknownIntradayInformationType || i >= InvalidIntradayInformationType {
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
	case Volume:
		return "5. volume"
	default:
		panic("invalid value of IntradayInformationType")
	}
}
