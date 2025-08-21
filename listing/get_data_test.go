package listing

import (
	"bytes"
	"os"
	"testing"
)

func TestGetData(t *testing.T) {

	data, err := os.ReadFile("../example_data/listing_status.csv")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var o = defaultOptions

	result, err := parseListingCsv(bytes.NewReader(data), &o)
	if err != nil {
		t.Fatalf("failed to parse CSV data: %v", err)
	}

	if result == nil {
		t.Fatal("result is unexpectedly nil")
	}
	if result.Meta == nil {
		t.Fatal("result.Meta is unexpectedly nil")
	}
	if len(result.Tradeables) == 0 {
		t.Fatal("result.Tradeables is unexpectedly empty")
	}
}

func TestGetData_1(t *testing.T) {

	data, err := os.ReadFile("../example_data/listing_status.csv")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var nyse = ExchangeName("NYSE")

	var o = defaultOptions
	o.ExchangeFilter = []ExchangeName{
		nyse,
	}

	result, err := parseListingCsv(bytes.NewReader(data), &o)
	if err != nil {
		t.Fatalf("failed to parse CSV data: %v", err)
	}

	if result == nil {
		t.Fatal("result is unexpectedly nil")
	}
	if result.Meta == nil {
		t.Fatal("result.Meta is unexpectedly nil")
	}
	if len(result.Tradeables) == 0 {
		t.Fatal("result.Tradeables is unexpectedly empty")
	}

	for _, v := range result.Tradeables {
		if v.Exchange != nyse {
			t.Fatalf("result.Tradeables has non-NYSE entries: %v", v)
		}
	}
}

func TestGetData_2(t *testing.T) {

	data, err := os.ReadFile("../example_data/listing_status.csv")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var stk = Stock

	var o = defaultOptions
	o.TypeFilter = []AssetType{stk}

	result, err := parseListingCsv(bytes.NewReader(data), &o)
	if err != nil {
		t.Fatalf("failed to parse CSV data: %v", err)
	}

	if result == nil {
		t.Fatal("result is unexpectedly nil")
	}
	if result.Meta == nil {
		t.Fatal("result.Meta is unexpectedly nil")
	}
	if len(result.Tradeables) == 0 {
		t.Fatal("result.Tradeables is unexpectedly empty")
	}

	for _, v := range result.Tradeables {
		if v.Type != stk {
			t.Fatalf("result.Tradeables has non-Stock entries: %v", v)
		}
	}
}
