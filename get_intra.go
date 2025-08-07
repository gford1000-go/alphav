package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/intraday"
)

func GetIntradayData(ctx context.Context, symbol string, opts ...func(*intraday.Options) error) (*intraday.Data, error) {

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return intraday.GetData(symbol, apiKey, opts...)

}
