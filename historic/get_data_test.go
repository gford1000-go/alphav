package historic

import (
	"os"
	"testing"

	"github.com/gford1000-go/alphav/common"
)

func TestParseJSON(t *testing.T) {

	data, err := os.ReadFile("../example_data/ibm_history.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	o := &Options{
		AllAvailableHistory: true,
	}

	result, err := parseJSON(data, o)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(result.TimeSeries) == 0 {
		t.Fatal("expected non-empty time series")
	}

	if result.Meta.Symbol != "IBM" {
		t.Fatalf("expected symbol 'IBM', got '%s'", result.Meta.Symbol)
	}

	if result.Meta.TimeZone != "US/Eastern" {
		t.Fatalf("expected timezone 'US/Eastern', got '%s'", result.Meta.TimeZone)
	}

	dt, _ := common.ParseDate("2025-08-19")
	if result.Meta.LastRefresh != dt {
		t.Fatalf("expected last refresh '2025-08-19', got '%s'", result.Meta.LastRefresh)
	}

	last, _ := common.ParseDate("2025-08-19")
	if result.Meta.DataRange.End != last {
		t.Fatalf("expected data range end '2025-08-19', got '%s'", result.Meta.DataRange.End)
	}

	first, _ := common.ParseDate("1999-11-01")
	if result.Meta.DataRange.Start != first {
		t.Fatalf("expected data range start '2025-01-01', got '%s'", result.Meta.DataRange.Start)
	}
}
