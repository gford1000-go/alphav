package alphav

import (
	"context"

	"github.com/gford1000-go/alphav/listing"
)

// GetActiveListing returns
func GetActiveListing(ctx context.Context, opts ...func(*listing.Options) error) (*listing.Data, error) {

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return listing.GetActiveListing(apiKey, opts...)
}
