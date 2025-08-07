package alphav

import (
	"context"
	"errors"
	"testing"
)

func TestInitialise(t *testing.T) {

	apiKey := "A KEY"

	ctx := Initialise(context.Background(), apiKey)

	apiKey1, err := getAPIKey(ctx)

	if err != nil {
		t.Fatalf("unexpected error: got: %v", err)
	}

	if apiKey != apiKey1 {
		t.Fatalf("mismatch in keys: expected: %s, got: %s", apiKey, apiKey1)
	}
}

func TestGetAPIKey(t *testing.T) {

	_, err := getAPIKey(context.Background())

	if err == nil {
		t.Fatalf("expected error: got nil")
	}

	if !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("unexpected error received.  expected: %v, got %v", ErrMissingAPIKey, err)
	}
}
