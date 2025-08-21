package listing

import "time"

// Metadata describes how the request was made
type Metadata struct {
	// Options describes the options used
	Options *Options
}

// Symbol is the type of the tradeable identifier
type Symbol string

// ExchangeName is the type of an exchange's name
type ExchangeName string

// Info is the set of information available for each symbol
type Info struct {
	Symbol   Symbol
	Name     string
	Exchange ExchangeName
	Type     AssetType
	IPO      time.Time
	Delisted time.Time
}

// Data is the returned object from a call to GetData
type Data struct {
	// Meta describes the details of the data
	Meta *Metadata
	// Tradeables is the map of all accessible data items
	Tradeables map[Symbol]*Info
}

func (d *Data) isValid() bool {
	return true
}
