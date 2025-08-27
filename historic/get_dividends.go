package historic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/gford1000-go/alphav/common"
)

// respDividendsJSON captures all possible return JSON
type respDividendsJSON struct {
	Info   *string         `json:"Information"`
	Err    *string         `json:"Error Message"`
	Meta   *metaJSON       `json:"Meta Data"`
	Symbol *string         `json:"Symbol"`
	Data   *[]*respDivJSON `json:"Data"`
}

type respDivJSON struct {
	ExDivDate       string `json:"ex_dividend_date"`
	DeclarationDate string `json:"declaration_date"`
	RecordDate      string `json:"record_date"`
	PaymentDate     string `json:"payment_date"`
	Amount          string `json:"amount"`
}

// GetDividends uses the provided apiKey to retrieve dividend details for the symbol
func GetDividends(symbol, apiKey string) (*DividendData, error) {

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=DIVIDENDS&symbol=%s&apikey=%s", symbol, apiKey)

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

	return parseDividendsJSON(b)
}

func parseDividendsJSON(b []byte) (*DividendData, error) {
	var d respDividendsJSON
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrParseError)
	}
	if d.Err != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Err, common.ErrRemoteCallError)
	}
	if d.Info != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Info, common.ErrRemoteCallError)
	}
	if d.Symbol == nil {
		return nil, fmt.Errorf("api error: expected Symbol, got nil: %w", common.ErrRemoteCallError)
	}
	if d.Data == nil || len(*d.Data) == 0 {
		return nil, fmt.Errorf("api error: expected Data, got empty list: %w", common.ErrRemoteCallError)
	}

	result := &DividendData{
		Meta: &Metadata{
			Symbol: *d.Symbol,
		},
		TimeSeries: []*DividendElement{},
	}

	if err := parseDividendTimeSeries(d.Data, result); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
	}

	return result, nil
}

func parseDividendTimeSeries(d *[]*respDivJSON, r *DividendData) error {

	if d == nil || len(*d) == 0 {
		return errors.New("no data available to be parsed")
	}

	ts := []*DividendElement{}

	dtRng := &DataRange{
		Start: latestDate,
		End:   earliestDate,
	}

	dtUpdate := earliestDate

	for i, v := range *d {
		if v == nil {
			return fmt.Errorf("v is nil for element %d", i)
		}

		ele := &DividendElement{}

		if v.RecordDate == "None" {
			ele.RecordDate = undefinedDate
		} else {
			t, err := common.ParseDate(v.RecordDate)
			if err != nil {
				return err
			}
			ele.RecordDate = DividendDate(t)

			if t.After(dtUpdate) {
				dtUpdate = t
			}
		}

		if v.DeclarationDate == "None" {
			ele.DeclarationDate = undefinedDate
		} else {
			t, err := common.ParseDate(v.DeclarationDate)
			if err != nil {
				return err
			}
			ele.DeclarationDate = DividendDate(t)
		}

		if v.ExDivDate == "None" {
			ele.ExDividendDate = undefinedDate
		} else {
			t, err := common.ParseDate(v.ExDivDate)
			if err != nil {
				return err
			}
			ele.ExDividendDate = DividendDate(t)

			if t.Before(dtRng.Start) {
				dtRng.Start = t
			}
			if t.After(dtRng.End) {
				dtRng.End = t
			}
		}

		if v.PaymentDate == "None" {
			ele.PaymentDate = undefinedDate
		} else {
			t, err := common.ParseDate(v.PaymentDate)
			if err != nil {
				return err
			}
			ele.PaymentDate = DividendDate(t)
		}

		value, err := strconv.ParseFloat(v.Amount, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount (%s) for element %d: %v", v.Amount, i, err)
		}
		ele.Amount = value

		ts = append(ts, ele)
	}

	// Sort is descending ... most recent date first, based on payment date
	slices.SortFunc(ts, func(a, b *DividendElement) int {
		return b.PaymentDate.Compare(a.PaymentDate)
	})

	r.TimeSeries = ts
	r.Meta.DataRange = dtRng      // Range based on ex-div dates
	r.Meta.LastRefresh = dtUpdate // LastRefresh based on latest RecordDate
	return nil
}
