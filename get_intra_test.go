package alphav

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gford1000-go/alphav/intraday"
)

func TestMain(m *testing.M) {
	// *** Add your Alpha Vantage TEST KEY here ***
	// Note that free api keys are limited to 25 requests/day.  This is based on IP address not the key itself.
	// Higher use requires premium access: see https://www.alphavantage.co/premium/
	os.Setenv("AV_API_KEY", "ATESTAPIKEY")

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
