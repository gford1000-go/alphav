package intraday

import "fmt"

type Interval int

const (
	UnknownInterval Interval = iota
	OneMin
	FiveMin
	FifteenMin
	ThirtyMin
	SixtyMin
	InvalidInterval
)

func (i Interval) String() string {
	switch i {
	case OneMin:
		return "1min"
	case FiveMin:
		return "5min"
	case FifteenMin:
		return "15min"
	case ThirtyMin:
		return "30min"
	case SixtyMin:
		return "60min"
	default:
		panic("invalid value of IntradayInterval")
	}
}

func parseInterval(s string) (Interval, error) {
	switch s {
	case "1min":
		return OneMin, nil
	case "5min":
		return FiveMin, nil
	case "15min":
		return FifteenMin, nil
	case "30min":
		return ThirtyMin, nil
	case "60min":
		return SixtyMin, nil
	default:
		return UnknownInterval, fmt.Errorf("unparseable interval: %s", s)
	}
}
