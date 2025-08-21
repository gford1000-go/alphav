package intraday

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

// metaJSON captures returned metadata, so it can be parsed
type metaJSON struct {
	Symbol   string `json:"2. Symbol"`
	Refresh  string `json:"3. Last Refreshed"`
	Interval string `json:"4. Interval"`
	Output   string `json:"5. Output Size"`
	TZ       string `json:"6. Time Zone"`
}

// respJson captures all possible return JSON
type respJSON struct {
	Info *string   `json:"Information"`
	Err  *string   `json:"Error Message"`
	Meta *metaJSON `json:"Meta Data"`
	TS1  any       `json:"Time Series (1min)"`
	TS5  any       `json:"Time Series (5min)"`
	TS15 any       `json:"Time Series (15min)"`
	TS30 any       `json:"Time Series (30min)"`
	TS60 any       `json:"Time Series (60min)"`
}

// GetData uses the provided apiKey to retrieve details for the symbol
func GetData(symbol, apiKey string, opts ...func(*Options) error) (*Data, error) {

	o := defaultOptions
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	outputsize := "compact"
	if o.RequestType {
		outputsize = "full"
	}

	month := ""
	if o.FromYear > 0 && o.FromMonth > 0 {
		month = fmt.Sprintf("&month=%d-%02d", o.FromYear, o.FromMonth)
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY%s&symbol=%s&interval=%s&adjusted=%v&extended_hours=%v&outputsize=%s&apikey=%s", month, symbol, o.Interval, o.Adjusted, o.ExtendedHours, outputsize, apiKey)

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

	var d respJSON
	if err = json.Unmarshal(b, &d); err != nil {
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

	if err := parseMetadata(d.Meta, result, &o); err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrMetadataParseError)
	}

	switch o.Interval {
	case OneMin:
		if err := parseTimeSeries(d.TS1, result, &o); err != nil {
			return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
		}
	case FiveMin:
		if err := parseTimeSeries(d.TS5, result, &o); err != nil {
			return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
		}
	case FifteenMin:
		if err := parseTimeSeries(d.TS15, result, &o); err != nil {
			return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
		}
	case ThirtyMin:
		if err := parseTimeSeries(d.TS30, result, &o); err != nil {
			return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
		}
	case SixtyMin:
		if err := parseTimeSeries(d.TS60, result, &o); err != nil {
			return nil, fmt.Errorf("%v: %w", err, common.ErrTimeSeriesParseError)
		}
	}

	return result, nil
}

func parseMetadata(m *metaJSON, r *Data, o *Options) error {
	im := &Metadata{
		Information: append([]InformationType{}, o.Information...),
		Symbol:      m.Symbol,
		TimeZone:    m.TZ,
	}

	interval, err := parseInterval(m.Interval)
	if err != nil {
		return err
	}
	im.RefreshInterval = interval

	im.LastRefresh, err = common.ParseIntradayDate(m.Refresh)
	if err != nil {
		return err
	}

	r.Meta = im
	return nil
}

func parseTimeSeries(i any, r *Data, o *Options) error {

	if i == nil {
		return errors.New("no data available to be parsed")
	}

	tmDataMap, ok := i.(map[string]any)
	if !ok {
		return errors.New("provided time series is of the wrong type")
	}

	tm := []*Element{}

	for k, v := range tmDataMap {
		if v == nil {
			return fmt.Errorf("v is nil for %s", k)
		}

		ele := &Element{
			Data: map[InformationType]float64{},
		}

		t, err := common.ParseIntradayDate(k)
		if err != nil {
			return err
		}
		ele.Timestamp = t

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
		return b.Timestamp.Compare(a.Timestamp)
	})

	r.TimeSeries = tm
	return nil
}
