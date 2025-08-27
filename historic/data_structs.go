package historic

import "time"

// Metadata describes what information was returned
type Metadata struct {
	// Information specifies the selected InformationTypes
	Information []InformationType
	// Symbol is the requested symbol for which data is retrieved
	Symbol string
	// LastRefresh is the time the data itself was last updated
	LastRefresh time.Time
	// TimeZone is the time zone of any returned datetime values
	TimeZone string
	// DataRange describes the range of data that was returned
	DataRange *DataRange
}

// Data Range describes the range of data history
type DataRange struct {
	// Start is the start time of the requested data range (earliest)
	Start time.Time
	// End is the end time of the requested data range (latest)
	End time.Time
}

// Element is an entry in the TimeSeries
type Element struct {
	// Timestamp is the time of the data in the element, based on the TZ in the metadata
	Date time.Time
	// Data holds the information for the specified types
	Data map[InformationType]float64
}

// Data is the returned object from a call to GetData
type Data struct {
	// Meta describes the details of the data
	Meta *Metadata
	// TimeSeries is an ordered set of data
	TimeSeries []*Element
}

func (d *Data) isValid() bool {
	if len(d.TimeSeries) == 0 || d.Meta == nil || len(d.Meta.Information) == 0 {
		return false
	}
	return true
}

// DividendDate allows for undefined date values to be identifiable in a series
type DividendDate time.Time

// IsUndefined returns true if the date was parseable from the timeseries details
func (d DividendDate) IsUndefined() bool {
	return d == undefinedDate
}

var undefinedDate = DividendDate(time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC))

// Compare returns:
//
//	1 if d is greater than other
//	0 if d is equal to other, or if either are undefined
//	-1 if d is less than other
func (d DividendDate) Compare(other DividendDate) int {
	if d.IsUndefined() || other.IsUndefined() {
		return 0
	}
	return time.Time(d).Compare(time.Time(other))
}

// DividendElement stores information related to a single dividend
type DividendElement struct {
	// RecordDate is the date that this information was stored
	RecordDate DividendDate
	// DeclarationDate is the date that the entity announced the dividend amount
	DeclarationDate DividendDate
	// ExDividendDate is the date on which ownership of the stock is determined, in terms of receipt of the amount
	ExDividendDate DividendDate
	// PaymentDate is the date on which the dividend is paid to the owner as identified on the ExDividendDate
	PaymentDate DividendDate
	// Amount is the amount of the dividend to be paid on the PaymentDate
	Amount float64
}

// DividendData is the history of dividend payments for the Symbol
type DividendData struct {
	// Meta describes the details of the data
	Meta *Metadata
	// TimeSeries is an ordered set of data
	TimeSeries []*DividendElement
}

func (d *DividendData) isValid() bool {
	if len(d.TimeSeries) == 0 || d.Meta == nil {
		return false
	}
	return true
}
