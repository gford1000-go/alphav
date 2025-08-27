package fx

import (
	"os"
	"testing"

	"github.com/gford1000-go/alphav/common"
)

func TestParseIntradayJSON(t *testing.T) {

	data, err := os.ReadFile("../example_data/usd_jpy_intraday.json")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	result, err := parseIntradayJSON(data)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(result.Data) == 0 {
		t.Fatal("expected non-empty data")
	}

	if result.Meta.FromCurrency != "USD" {
		t.Fatalf("expected symbol 'USD', got '%s'", result.Meta.FromCurrency)
	}

	if result.Meta.ToCurrency != "JPY" {
		t.Fatalf("expected symbol 'JPY', got '%s'", result.Meta.ToCurrency)
	}

	if result.Meta.TimeZone != "UTC" {
		t.Fatalf("expected timezone 'UTC', got '%s'", result.Meta.TimeZone)
	}

	dt, _ := common.ParseIntradayDate("2025-08-27 11:25:19")
	if result.Meta.LastRefresh != dt {
		t.Fatalf("expected last refresh '2025-08-27 11:25:19', got '%s'", result.Meta.LastRefresh)
	}

	if len(result.Data) == 0 {
		t.Fatal("expected data, got zero length")
	}

	bid, ask, rate := 148.103, 148.11, 148.11
	if !common.EqualFloat64(bid, result.Data[Bid], 4) {
		t.Fatalf("expected data element to have bid: '%v', got '%v'", bid, result.Data[Bid])
	}
	if !common.EqualFloat64(ask, result.Data[Ask], 4) {
		t.Fatalf("expected first data element to have high: '%v', got '%v'", ask, result.Data[Ask])
	}
	if !common.EqualFloat64(rate, result.Data[FXRate], 4) {
		t.Fatalf("expected first data element to have low: '%v', got '%v'", rate, result.Data[FXRate])
	}
}
