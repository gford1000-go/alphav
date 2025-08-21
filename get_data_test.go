package alphav

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gford1000-go/alphav/historic"
	"github.com/gford1000-go/alphav/intraday"
	"github.com/gford1000-go/alphav/listing"
)

func TestMain(m *testing.M) {
	// *** Add your Alpha Vantage TEST KEY here ***
	// Note that free api keys are limited to 25 requests/day.  This is based on IP address not the key itself.
	// Higher use requires premium access: see https://www.alphavantage.co/premium/
	os.Setenv("AV_API_KEY", "ADD_TEST_KEY_HERE")

	// Run tests
	code := m.Run()

	// Clean up if needed
	os.Unsetenv("AV_API_KEY")

	os.Exit(code)
}

func ExampleGetIntradayData() {

	apiKey := os.Getenv("AV_API_KEY")

	ctx := Initialise(context.Background(), apiKey)

	if data, err := GetIntradayData(ctx, "IBM", intraday.WithExtendedHours(false)); err == nil {
		fmt.Println(len(data.TimeSeries))
	} else {
		fmt.Println(err)
	}

	// Output:
	// 100
}

func ExampleGetHistoricData() {

	apiKey := os.Getenv("AV_API_KEY")

	ctx := Initialise(context.Background(), apiKey)

	if data, err := GetHistoricData(ctx, "IBM", historic.WithAllAvailableHistory(false)); err == nil {
		fmt.Println(len(data.TimeSeries))
	} else {
		fmt.Println(err)
	}

	// Output:
	// 100
}

func ExampleGetActiveListing() {

	apiKey := os.Getenv("AV_API_KEY")

	ctx := Initialise(context.Background(), apiKey)

	if data, err := GetActiveListing(ctx, listing.WithOnlyTypes([]listing.AssetType{listing.ETF})); err == nil {
		fmt.Println(len(data.Tradeables) > 0)
	} else {
		fmt.Println(err)
	}

	// Output:
	// true
}
