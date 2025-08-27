package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/fx"
)

// GetIntradayData returns data for the specified currency pair, using the api_key stored in the context.
// opts allows the behaviour of the call to be varied per the options in https://www.alphavantage.co/documentation/
// for FX_DAILY
func GetFX(ctx context.Context, fromCurrency, toCurrency string, opts ...func(*fx.Options) error) (*fx.Data, error) {

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return fx.GetData(fromCurrency, toCurrency, apiKey, opts...)

}
