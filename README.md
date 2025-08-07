[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/gford1000-go/alphav)

# alphav

Provides a wrapper to [Alpha Vantage](https://www.alphavantage.co) api, specifically the `TIME_SERIES_INTRADAY` request.

Note: an API KEY is required to use the API, with free keys rate limited to 25 request/day.

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
