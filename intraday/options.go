package intraday

import (
	"fmt"
	"time"

	"github.com/gford1000-go/alphav/common"
)

// Options can change the returned Data from GetData
type Options struct {
	// Interval specifies the interval between time series elements.  Default: 5min
	Interval Interval
	// Adjusted = true indicates split/dividend adjusted data is returned.  Default: true
	Adjusted bool
	// ExtendedHours = true indicates include close market activity.  Default: false
	ExtendedHours bool
	// Information specifies the set of data types to be returned.  Default: all types
	Information []InformationType
	// RequestType = true returns 30 days worth of data; false is 100 records.  Default: false
	RequestType bool
	// FromYear, if set, specifies the start year from which data is to be returned.  Default: current year
	FromYear int
	// FromMonth, if set, specifies the start month from wihc data is to be returned.  Default: latest data
	FromMonth int
}

// WithInterval sets the interval between elements
func WithInterval(interval Interval) func(*Options) error {
	return func(o *Options) error {
		if !interval.isValid() {
			return common.ErrInvalidInterval
		}
		o.Interval = interval
		return nil
	}
}

// WithAdjusted specifies adjusted when true
func WithAdjusted(adjusted bool) func(*Options) error {
	return func(o *Options) error {
		o.Adjusted = adjusted
		return nil
	}
}

// WithExtendedHours specifies to include closed market details.
func WithExtendedHours(extended bool) func(*Options) error {
	return func(o *Options) error {
		o.ExtendedHours = extended
		return nil
	}
}

func WithRequestType(requestType bool) func(*Options) error {
	return func(o *Options) error {
		o.RequestType = requestType
		return nil
	}
}

// WithStartPoint sets the YYYY-MM point from which data is retrieved.  This can be any month after 2000-01.
func WithStartPoint(year, month int) func(*Options) error {
	return func(o *Options) error {
		if year > time.Now().Year() || year < 2000 {
			return fmt.Errorf("invalid start year specified: %d", year)
		}
		if month < 1 || month > 12 {
			return fmt.Errorf("invalid start month specified: %d", month)
		}
		o.FromMonth = month
		o.FromYear = year
		return nil
	}
}

// WithInformation sets the information types to be returned
// If no information types are specified, all types are returned.
func WithInformation(information ...InformationType) func(*Options) error {
	return func(o *Options) error {
		for _, i := range information {
			if !i.isValid() {
				return common.ErrInvalidInformationType
			}
		}
		o.Information = append([]InformationType{}, o.Information...)
		return nil
	}
}

var defaultOptions = Options{
	Interval:      FiveMin,
	Adjusted:      true,
	ExtendedHours: false,
	Information: []InformationType{
		Open,
		High,
		Low,
		Close,
		Volume,
	},
	RequestType: false,
}
