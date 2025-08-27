package fx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gford1000-go/alphav/common"
)

// metaJSON captures returned metadata, so it can be parsed
type metaJSON struct {
	From    string `json:"2. From Symbol"`
	To      string `json:"3. To Symbol"`
	Refresh string `json:"5. Last Refreshed"`
	Output  string `json:"4. Output Size"`
	TZ      string `json:"6. Time Zone"`
}

// respJson captures all possible return JSON
type respJSON struct {
	Info *string   `json:"Information"`
	Err  *string   `json:"Error Message"`
	Meta *metaJSON `json:"Meta Data"`
	TSD  any       `json:"Time Series FX (Daily)"`
}

// GetData uses the provided apiKey to retrieve details for the symbol
func GetData(fromCurrency, toCurrency, apiKey string, opts ...func(*Options) error) (*Data, error) {

	o := defaultOptions
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	outputsize := "compact"
	if o.AllAvailableHistory {
		outputsize = "full"
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=%s&to_symbol=%s&outputsize=%s&apikey=%s", strings.ToUpper(fromCurrency), strings.ToUpper(toCurrency), outputsize, apiKey)

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

	return parseJSON(b, &o)
}

func parseJSON(b []byte, o *Options) (*Data, error) {
	var d respJSON
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrParseError)
	}
	if d.Err != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Err, common.ErrRemoteCallError)
	}
	if d.Info != nil {
		return nil, fmt.Errorf("api error: %s: %w", *d.Info, common.ErrRemoteCallError)
	}

	result := &Data{
		Meta:       &Metadata{},
		TimeSeries: []*Element{},
	}

	if err := parseMetadata(d.Meta, result, o); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrMetadataParseError)
	}

	if err := parseTimeSeries(d.TSD, result, o); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
	}

	return result, nil
}

func parseMetadata(m *metaJSON, r *Data, o *Options) error {
	im := &Metadata{
		Information:  append([]InformationType{}, o.Information...),
		FromCurrency: m.From,
		ToCurrency:   m.To,
		TimeZone:     m.TZ,
	}

	var err error
	im.LastRefresh, err = common.ParseDate(m.Refresh)
	if err != nil {
		return err
	}

	r.Meta = im
	return nil
}

var earliestDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
var latestDate = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

func parseTimeSeries(i any, r *Data, o *Options) error {

	if i == nil {
		return errors.New("no data available to be parsed")
	}

	tmDataMap, ok := i.(map[string]any)
	if !ok {
		return errors.New("provided time series is of the wrong type")
	}

	tm := []*Element{}

	dtRng := &DataRange{
		Start: latestDate,
		End:   earliestDate,
	}

	for k, v := range tmDataMap {
		if v == nil {
			return fmt.Errorf("v is nil for %s", k)
		}

		ele := &Element{
			Data: map[InformationType]float64{},
		}

		t, err := common.ParseDate(k)
		if err != nil {
			return err
		}
		ele.Date = t

		if t.Before(dtRng.Start) {
			dtRng.Start = t
		}
		if t.After(dtRng.End) {
			dtRng.End = t
		}

		m, ok := v.(map[string]any)
		if !ok {
			return errors.New("failed to extract data from data map")
		}

		for _, it := range o.Information {
			mv, ok := m[it.toAVString()]
			if !ok {
				return fmt.Errorf("missing %s for %s", it, k)
			}

			s, ok := mv.(string)
			if !ok {
				return fmt.Errorf("value of %s for %s is not a string type (%v)", it, k, mv)
			}

			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("error parsing %s (%s) for %s: %v", mv, it, k, err)
			}
			ele.Data[it] = value
		}

		tm = append(tm, ele)
	}

	// Sort is descending ... most recent date first
	slices.SortFunc(tm, func(a, b *Element) int {
		return b.Date.Compare(a.Date)
	})

	r.TimeSeries = tm
	r.Meta.DataRange = dtRng
	return nil
}
