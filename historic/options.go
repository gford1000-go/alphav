package historic

// Options can change the returned Data from GetData
type Options struct {
	// Information specifies the set of data types to be returned.  Default: all types
	Information []InformationType
	// AllAvailableHistory = true returns 20 years worth of data; false is 100 records.  Default: false
	AllAvailableHistory bool
}

func WithAllAvailableHistory(all bool) func(*Options) error {
	return func(o *Options) error {
		o.AllAvailableHistory = all
		return nil
	}
}

var defaultOptions = Options{
	Information: []InformationType{
		Open,
		High,
		Low,
		Close,
		Volume,
		AdjustedClose,
		DividendAmount,
		SplitCoefficient,
	},
	AllAvailableHistory: false,
}
