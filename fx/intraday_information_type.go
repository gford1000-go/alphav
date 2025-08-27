package fx

type IntradayInformationType int

const (
	UnknownIntradayInformationType IntradayInformationType = iota
	FXRate
	Bid
	Ask
	InvalidIntradayInformationType
)

func (i IntradayInformationType) String() string {
	switch i {
	case FXRate:
		return "rate"
	case Bid:
		return "bid"
	case Ask:
		return "ask"
	default:
		panic("invalid value of IntradayInformationType")
	}
}

func (i IntradayInformationType) isValid() bool {
	if i <= UnknownIntradayInformationType || i >= InvalidIntradayInformationType {
		return false
	}
	return true
}
