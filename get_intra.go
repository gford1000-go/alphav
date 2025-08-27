package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/common"
	"github.com/gford1000-go/alphav/intraday"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetIntradayData returns data for the specified symbol, using the api_key stored in the context.
// opts allows the behaviour of the call to be varied per the options in https://www.alphavantage.co/documentation/
// for TIME_SERIES_INTRADAY
func GetIntradayData(ctx context.Context, symbol string, opts ...func(*intraday.Options) error) (*intraday.Data, error) {

	tracer := otel.Tracer(common.TracerName)

	ctx, span := tracer.Start(ctx, "GetIntradayData")
	defer span.End()

	span.SetAttributes(attribute.String("Symbol", symbol))

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return intraday.GetData(symbol, apiKey, opts...)

}
