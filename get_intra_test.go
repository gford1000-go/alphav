package alphav

import (
	"context"
	"fmt"
	"os"

	"github.com/gford1000-go/alphav/intraday"
)

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
