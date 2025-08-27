package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/common"
	"github.com/gford1000-go/alphav/listing"
	"go.opentelemetry.io/otel"
)

// GetActiveListing returns
func GetActiveListing(ctx context.Context, opts ...func(*listing.Options) error) (*listing.Data, error) {

	tracer := otel.Tracer(common.TracerName)

	ctx, span := tracer.Start(ctx, "GetActiveListing")
	defer span.End()

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return listing.GetActiveListing(apiKey, opts...)
}
