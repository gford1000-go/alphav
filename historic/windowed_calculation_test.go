package historic

import (
	"context"
	"math"
	"os"
	"testing"
	"time"
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

func TestGetWindowedCalculation_1(t *testing.T) {

	history, _ := os.ReadFile("../example_data/ibm_history.json")

	var o Options = defaultOptions
	o.AllAvailableHistory = true

	data, err := parseJSON(history, &o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var tag = "Change"
	results, err := GetWindowedCalculation(context.Background(), data, 1, AdjustedClose, map[string]WindowFunc{
		tag: WindowChange,
	},
		WithNumberOfDataPoints(1))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results.TimeSeries[tag]) != 1 {
		t.Fatalf("unexpected length of time series returned: expected 1, got %d", len(results.TimeSeries[tag]))
	}

	ele := results.TimeSeries[tag][0]
	if ele.WindowStart != time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected date result: wanted: 2025-08-19, got: %v", ele.WindowStart)
	}

	var roundedChange = math.Round(ele.Value*100) / 100
	if roundedChange != 1.83 {
		t.Fatalf("unexpected value result: wanted: 1.83, got: %v", roundedChange)
	}
}
