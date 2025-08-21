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
		WithElementProcessingLimit(1))

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

func TestGetWindowedCalculation_2(t *testing.T) {

	history, _ := os.ReadFile("../example_data/ibm_history.json")

	var o Options = defaultOptions
	o.AllAvailableHistory = true

	data, err := parseJSON(history, &o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var windowLength = 1      // Window length of 1 Element (i.e. day on day change)
	var limitOfDataPoints = 1 // Number of data elements to be processed from most recent
	var tag1 = "Change"
	var tag2 = "% Change"
	results, err := GetWindowedCalculation(context.Background(), data, windowLength, AdjustedClose, map[string]WindowFunc{
		tag1: WindowChange,
		tag2: WindowPercentageChange},
		WithElementProcessingLimit(limitOfDataPoints))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results.TimeSeries[tag1]) != limitOfDataPoints {
		t.Fatalf("unexpected length of time series returned: expected %d, got %d", limitOfDataPoints, len(results.TimeSeries[tag1]))
	}

	ele := results.TimeSeries[tag1][0]
	if ele.WindowStart != time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected date result: wanted: 2025-08-19, got: %v", ele.WindowStart)
	}

	var roundedChange = math.Round(ele.Value*100) / 100
	if roundedChange != 1.83 {
		t.Fatalf("unexpected value result: wanted: 1.83, got: %v", roundedChange)
	}

	if len(results.TimeSeries[tag2]) != limitOfDataPoints {
		t.Fatalf("unexpected length of time series returned: expected %d, got %d", limitOfDataPoints, len(results.TimeSeries[tag2]))
	}

	ele2 := results.TimeSeries[tag2][0]
	if ele2.WindowStart != time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected date result: wanted: 2025-08-19, got: %v", ele2.WindowStart)
	}

	var roundedPercentageChange = math.Round(ele2.Value*100) / 100
	if roundedPercentageChange != 0.76 {
		t.Fatalf("unexpected value result: wanted: 0.76, got: %v", roundedPercentageChange)
	}

}

func TestGetWindowedCalculation_3(t *testing.T) {

	history, _ := os.ReadFile("../example_data/ibm_history.json")

	var o Options = defaultOptions
	o.AllAvailableHistory = true

	data, err := parseJSON(history, &o)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var windowLength = 15  // Change over 15 Elements
	var numDataPoints = 30 // Limit to processing over first 30 elements of the time series
	var tag1 = "Change"
	var tag2 = "% Change"
	results, err := GetWindowedCalculation(context.Background(), data, windowLength, AdjustedClose, map[string]WindowFunc{
		tag1: WindowChange,
		tag2: WindowPercentageChange},
		WithElementProcessingLimit(numDataPoints))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results.TimeSeries[tag1]) != numDataPoints {
		t.Fatalf("unexpected length of time series returned: expected %d, got %d", numDataPoints, len(results.TimeSeries[tag1]))
	}

	ele := results.TimeSeries[tag1][0]
	if ele.WindowStart != time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected date result: wanted: 2025-08-19, got: %v", ele.WindowStart)
	}

	var roundedChange = math.Round(ele.Value*100) / 100
	if roundedChange != -19.32 {
		t.Fatalf("unexpected value result: wanted: -19.32, got: %v", roundedChange)
	}

	if len(results.TimeSeries[tag2]) != numDataPoints {
		t.Fatalf("unexpected length of time series returned: expected %d, got %d", numDataPoints, len(results.TimeSeries[tag2]))
	}

	ele2 := results.TimeSeries[tag2][0]
	if ele2.WindowStart != time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("unexpected date result: wanted: 2025-08-19, got: %v", ele2.WindowStart)
	}

	var roundedPercentageChange = math.Round(ele2.Value*100) / 100
	if roundedPercentageChange != -7.41 {
		t.Fatalf("unexpected value result: wanted: -7.41, got: %v", roundedPercentageChange)
	}

}
