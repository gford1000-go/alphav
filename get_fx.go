package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/common"
	"github.com/gford1000-go/alphav/fx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetFX returns data for the specified currency pair, using the api_key stored in the context.
// opts allows the behaviour of the call to be varied per the options in https://www.alphavantage.co/documentation/
// for FX_DAILY
func GetFX(ctx context.Context, fromCurrency, toCurrency string, opts ...func(*fx.Options) error) (*fx.Data, error) {

	tracer := otel.Tracer(common.TracerName)

	ctx, span := tracer.Start(ctx, "GetFX")
	defer span.End()

	span.SetAttributes(attribute.String("FromCurrency", fromCurrency))
	span.SetAttributes(attribute.String("ToCurrency", toCurrency))

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return fx.GetData(fromCurrency, toCurrency, apiKey, opts...)

}

// GetIntradayFX returns data for the specified currency pair, using the api_key stored in the context.
// This uses CURRENCY_EXCHANGE_RATE from https://www.alphavantage.co/documentation/
func GetIntradayFX(ctx context.Context, fromCurrency, toCurrency string) (*fx.IntradayData, error) {

	tracer := otel.Tracer(common.TracerName)

	ctx, span := tracer.Start(ctx, "GetIntradayFX")
	defer span.End()

	span.SetAttributes(attribute.String("FromCurrency", fromCurrency))
	span.SetAttributes(attribute.String("ToCurrency", toCurrency))

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return fx.GetIntraday(fromCurrency, toCurrency, apiKey)

}
