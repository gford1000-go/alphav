package listing

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gford1000-go/alphav/common"
)

func GetActiveListing(apiKey string, opts ...func(*Options) error) (*Data, error) {

	var o = defaultOptions
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=LISTING_STATUS&apikey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, common.ErrRemoteCallError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", resp.Status, common.ErrRemoteCallError)
	}

	return parseListingCsv(resp.Body, &o)
}

func parseListingCsv(data io.Reader, o *Options) (*Data, error) {

	reader := csv.NewReader(data)
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w: %w", err, common.ErrRemoteCallError)
	}

	colIndex := map[string]int{}
	for i, col := range header {
		colIndex[strings.ToLower(col)] = i
	}

	var tradeables = map[Symbol]*Info{}

	var line = 0
	for {
		line++
		record, err := reader.Read()
		if err == io.EOF {
			break // Done
		}
		if err != nil {
			return nil, err
		}

		if len(record) != len(header) {
			continue
		}

		if record[colIndex["status"]] != "Active" {
			continue // Must be active
		}

		ipo, err := common.ParseDate(record[colIndex["ipodate"]])
		if err != nil {
			return nil, fmt.Errorf("line: %d: ipoDate error parsing '%s': %v", line, record[colIndex["ipodate"]], err)
		}

		var delist time.Time
		delisted := record[colIndex["delistingdate"]]
		if delisted != "null" {
			delist, err = common.ParseDate(delisted)
			if err != nil {
				return nil, fmt.Errorf("line: %d: delistingDate error parsing '%s': %v", line, record[colIndex["delistingdate"]], err)
			}
		}

		assetType, err := parseAssetType(record[colIndex["assettype"]])
		if err != nil {
			return nil, fmt.Errorf("line: %d: assetType error parsing '%s': %v", line, record[colIndex["assettype"]], err)
		}

		var canAdd bool = true
		var exchange ExchangeName = ExchangeName(strings.ToUpper(record[colIndex["exchange"]]))
		if len(o.ExchangeFilter) > 0 {
			canAdd = false
			for _, name := range o.ExchangeFilter {
				if name == exchange {
					canAdd = true
					break
				}
			}
		}

		if !canAdd {
			continue // Filtered out by exchange
		}

		if len(o.TypeFilter) > 0 {
			canAdd = false
			for _, t := range o.TypeFilter {
				if t == assetType {
					canAdd = true
					break
				}
			}
		}

		if canAdd {
			info := &Info{
				Symbol:   Symbol(record[colIndex["symbol"]]),
				Name:     record[colIndex["name"]],
				Type:     assetType,
				Exchange: exchange,
				IPO:      ipo,
				Delisted: delist,
			}

			tradeables[info.Symbol] = info
		}
	}

	return &Data{
		Meta: &Metadata{
			Options: o,
		},
		Tradeables: tradeables,
	}, nil
}
