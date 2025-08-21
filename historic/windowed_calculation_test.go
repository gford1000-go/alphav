package historic

import (
	"context"
	"os"
	"testing"
)

func TestGetWindowedCalculation(t *testing.T) {

	history, _ := os.ReadFile("../example_data/ibm_history.json")

	var o Options = defaultOptions
	o.AllAvailableHistory = true

	data, err := parseJSON(history, &o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var avgTag = "Average"
	results, err := GetWindowedCalculation(context.Background(), data, 1, AdjustedClose, map[string]WindowFunc{
		avgTag: WindowAverage,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results.TimeSeries[avgTag]) != len(data.TimeSeries)-1 {
		t.Fatalf("unexpected length of time series returned: expected %d, got %d", len(data.TimeSeries)-1, len(results.TimeSeries[avgTag]))
	}
}
