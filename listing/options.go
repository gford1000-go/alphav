package listing

import (
	"errors"
	"strings"
)

// Options can change the returned Data from GetData
type Options struct {
	// TypeFilter limits to only the specified AssetTypes.  Default is all AssetTypes
	TypeFilter []AssetType
	// ExchangeName limits to only the specified Exchanges.  Default is all Exchanges
	ExchangeFilter []ExchangeName
}

var defaultOptions = Options{}

// WithOnlyTypes limits the set of returned listings to be restricted to the specified AssetTypes
// Not setting a type filter means all listings of any type are returned.
func WithOnlyTypes(types []AssetType) func(*Options) error {
	return func(o *Options) error {
		for _, t := range types {
			if !t.isValid() {
				return errors.New("invalid asset type")
			}
		}
		o.TypeFilter = types
		return nil
	}
}

// WithOnlyNamedExchanges limits the set of returned listings to those available from the
// specified exchanges.
// Not setting an exchange filter means all listings from all exchanges are returned
func WithOnlyNamedExchanges(exchanges []ExchangeName) func(*Options) error {
	return func(o *Options) error {
		if len(exchanges) > 0 {
			o.ExchangeFilter = []ExchangeName{}
			for _, exchange := range exchanges {
				o.ExchangeFilter = append(o.ExchangeFilter, ExchangeName(strings.ToUpper(string(exchange))))
			}
		}
		return nil
	}
}
