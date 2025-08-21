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
