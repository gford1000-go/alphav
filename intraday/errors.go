package intraday

import "errors"

// ErrRemoteCallError returned when errors are raised calling Alpha Vantage URL
var ErrRemoteCallError = errors.New("error calling URL")

// ErrParseError returned when JSON parsing errors occur
var ErrParseError = errors.New("unexpected data returned")

// ErrMetadataParseError returned when meta data parsing has errors
var ErrMetadataParseError = errors.New("invalid metadata received")

// ErrTimeSeriesParseError returned when time series parsing has errors
var ErrTimeSeriesParseError = errors.New("invalid time series received")
