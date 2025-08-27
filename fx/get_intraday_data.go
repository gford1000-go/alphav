package fx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gford1000-go/alphav/common"
)

type respIntradayJSON struct {
	Info *string       `json:"Information"`
	Err  *string       `json:"Error Message"`
	Data *intradayJSON `json:"Realtime Currency Exchange Rate"`
}

type intradayJSON struct {
	From     string `json:"1. From_Currency Code"`
	FromName string `json:"2. From_Currency Name"`
	To       string `json:"3. To_Currency Code"`
	ToName   string `json:"4. To_Currency Name"`
	Rate     string `json:"5. Exchange Rate"`
	Bid      string `json:"8. Bid Price"`
	Ask      string `json:"9. Ask Price"`
	Refresh  string `json:"6. Last Refreshed"`
	TZ       string `json:"7. Time Zone"`
}

// GetIntraday uses the provided apiKey to retrieve details for the symbol
func GetIntraday(fromCurrency, toCurrency, apiKey string) (*IntradayData, error) {

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_symbol=%s&to_symbol=%s&apikey=%s", strings.ToUpper(fromCurrency), strings.ToUpper(toCurrency), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrRemoteCallError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", resp.Status, common.ErrRemoteCallError)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrRemoteCallError)
	}

	return parseIntradayJSON(b)
}

func parseIntradayJSON(b []byte) (*IntradayData, error) {
	var d respIntradayJSON
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrParseError)
	}
	if d.Err != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Err, common.ErrRemoteCallError)
	}
	if d.Info != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Info, common.ErrRemoteCallError)
	}
	if d.Data == nil {
		return nil, errors.New("no data available to be parsed")
	}
	t, err := common.ParseIntradayDate(d.Data.Refresh)
	if err != nil {
		return nil, err
	}

	result := &IntradayData{
		Meta: &Metadata{
			FromCurrency: d.Data.From,
			ToCurrency:   d.Data.To,
			LastRefresh:  t,
			TimeZone:     d.Data.TZ,
			DataRange: &DataRange{
				Start: t,
				End:   t,
			},
		},
		Data: map[IntradayInformationType]float64{},
	}

	for _, t := range []IntradayInformationType{Bid, Ask, FXRate} {
		var v string
		switch t {
		case Bid:
			v = d.Data.Bid
		case Ask:
			v = d.Data.Ask
		case FXRate:
			v = d.Data.Rate
		default:
			return nil, errors.New("invalid value of IntradayInformationType requested")
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		result.Data[t] = f
	}

	return result, nil
}
