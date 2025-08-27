package historic

import (
	"os"
	"testing"

	"github.com/gford1000-go/alphav/common"
)

func TestParseDividendsJSON(t *testing.T) {

	data, err := os.ReadFile("../example_data/ibm_dividends.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	result, err := parseDividendsJSON(data)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(result.TimeSeries) == 0 {
		t.Fatal("expected non-empty time series")
	}

	if result.Meta.Symbol != "IBM" {
		t.Fatalf("expected symbol 'IBM', got '%s'", result.Meta.Symbol)
	}

	dt, _ := common.ParseDate("2025-08-08")
	if result.Meta.LastRefresh != dt {
		t.Fatalf("expected last refresh '2025-08-08', got '%s'", result.Meta.LastRefresh)
	}

	last, _ := common.ParseDate("2025-08-08")
	if result.Meta.DataRange.End != last {
		t.Fatalf("expected data range end '2025-08-08', got '%s'", result.Meta.DataRange.End)
	}

	first, _ := common.ParseDate("1999-02-08")
	if result.Meta.DataRange.Start != first {
		t.Fatalf("expected data range start '1999-02-08', got '%s'", result.Meta.DataRange.Start)
	}

	latest := 1.68
	if !common.EqualFloat64(latest, result.TimeSeries[0].Amount, 4) {
		t.Fatalf("expected latest data element to have amount: '%v', got '%v'", latest, result.TimeSeries[0].Amount)
	}

	earliest := 0.22
	if !common.EqualFloat64(earliest, result.TimeSeries[len(result.TimeSeries)-1].Amount, 4) {
		t.Fatalf("expected first data element to have amount: '%v', got '%v'", earliest, result.TimeSeries[len(result.TimeSeries)-1].Amount)
	}
}
