[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/gford1000-go/alphav)

# alphav

Provides a wrapper to the [Alpha Vantage](https://www.alphavantage.co) api.

The following are supported:

* `LISTING_STATUS`
* `TIME_SERIES_INTRADAY`
* `TIME_SERIES_DAILY_ADJUSTED` (requires a premium account)
* `FX_DAILY`
* `CURRENCY_EXCHANGE_RATE`
* `DIVIDENDS`

This allows the set of available tradeables to be retrieved, together with 20 year histories and recent intraday activity.

Note: an API KEY is required to use the API, with free keys rate limited to 25 request/day.

The package maintains a consistent behaviour across calls, for example:

```go
func main() {
    apiKey := "MY API KEY" // Provided by some means

    ctx := alphav.Initialise(context.Background(), apiKey)

    data, err := alphav.GetIntradayData(ctx, "IBM", intraday.WithExtendedHours(false))
    if err != nil {
        fmt.Println(err)
    }

    // Do something with data
}
```

See examples and tests for more details.
