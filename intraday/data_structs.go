package intraday

import "time"

// Metadata describes what information was returned
type Metadata struct {
	// Information specifies the selected InformationTypes
	Information []InformationType
	// Symbol is the requested symbol for which data is retrieved
	Symbol string
	// LastRefresh is the time the data itself was last updated
	LastRefresh time.Time
	// RefreshInterval specifies the selected Interval in the time series history
	RefreshInterval Interval
	// TimeZone is the time zone of any returned datetime values
	TimeZone string
}

// Element is an entry in the TimeSeries
type Element struct {
	// Timestamp is the time of the data in the element, based on the TZ in the metadata
	Timestamp time.Time
	// ExtendedHours is true if the element is out of main trading hours
	ExtendedHours bool
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
