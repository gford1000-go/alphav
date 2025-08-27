package fx

import (
	"os"
	"testing"
	"time"

	"github.com/gford1000-go/alphav/common"
)

func TestParseJSON(t *testing.T) {

	data, err := os.ReadFile("../example_data/eur_usd.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	o := &Options{
		AllAvailableHistory: true,
		Information: []InformationType{
			Open,
			Close,
			High,
			Low,
		},
	}

	result, err := parseJSON(data, o)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(result.TimeSeries) == 0 {
		t.Fatal("expected non-empty time series")
	}

	if result.Meta.FromCurrency != "EUR" {
		t.Fatalf("expected symbol 'EUR', got '%s'", result.Meta.FromCurrency)
	}

	if result.Meta.TimeZone != "UTC" {
		t.Fatalf("expected timezone 'UTC', got '%s'", result.Meta.TimeZone)
	}

	dt, _ := common.ParseDate("2025-08-26")
	if result.Meta.LastRefresh != dt {
		t.Fatalf("expected last refresh '2025-08-26', got '%s'", result.Meta.LastRefresh)
	}

	last, _ := common.ParseDate("2025-08-26")
	if result.Meta.DataRange.End != last {
		t.Fatalf("expected data range end '2025-08-26', got '%s'", result.Meta.DataRange.End)
	}

	first, _ := common.ParseDate("2025-04-09")
	if result.Meta.DataRange.Start != first {
		t.Fatalf("expected data range start '2025-04-09', got '%s'", result.Meta.DataRange.Start)
	}

	if len(result.TimeSeries) == 0 {
		t.Fatal("expected data, got zero length timeseries")
	}

	data26 := result.TimeSeries[0]
	if data26.Date != time.Date(2025, 8, 26, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("expected first data element to be '2025-08-26', got '%s'", data26.Date)
	}

	open, high, low, close := 1.1618, 1.1665, 1.1599, 1.1642
	if !common.EqualFloat64(open, data26.Data[Open], 4) {
		t.Fatalf("expected first data element to have open: '%v', got '%v'", open, data26.Data[Open])
	}
	if !common.EqualFloat64(high, data26.Data[High], 4) {
		t.Fatalf("expected first data element to have high: '%v', got '%v'", high, data26.Data[High])
	}
	if !common.EqualFloat64(low, data26.Data[Low], 4) {
		t.Fatalf("expected first data element to have low: '%v', got '%v'", low, data26.Data[Low])
	}
	if !common.EqualFloat64(close, data26.Data[Close], 4) {
		t.Fatalf("expected first data element to have close: '%v', got '%v'", close, data26.Data[Close])
	}

}
