package alphav

import (
	"context"
	"errors"
)

type apiKeyKey string

var apiKeyKeyName apiKeyKey = "theKey"

// Initialise registers the supplied API Key to Alpha Vantage
func Initialise(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, apiKeyKeyName, apiKey)
}

// ErrMissingAPIKey returned if the api key has not been found (Initialise() not called)
var ErrMissingAPIKey = errors.New("context did not contain a valid api key")

// getAPIKey retrieves the api key from the context
func getAPIKey(ctx context.Context) (string, error) {
	v := ctx.Value(apiKeyKeyName)
	if v != nil {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}
	return "", ErrMissingAPIKey
}
