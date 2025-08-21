package historic

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/gford1000-go/alphav/common"
)

// WindowsMeta describes what was passed to GetWindowedCalculation
type WindowedMeta struct {
	// Input is the time series data provided
	Input *Data
	// WindowLength is the interval used for calculations
	WindowLength int
	// InformationType is the value to be extracted from the time series for calculations
	InformationType InformationType
	// Calculations is the map of a tag value to the WindowFunc to generate the result
	Calculations map[string]WindowFunc
	// Options is the set of options used within the calculation processing
	Options *WindowedCalculationOptions
}

// WindowedElement is a single result from the time series processing
type WindowedElement struct {
	// WindowStart is the point at which calculations are being performed (i.e. the latest date)
	WindowStart time.Time
	// Value is the result of the calculation as of WindowStart
	Value float64
}

// WindowedResult is returned by GetWindowedCalculation
type WindowedResult struct {
	// Meta describes the input to the call
	Meta *WindowedMeta
	// TimeSeries is the set of results for each of the WindowFunc used
	TimeSeries map[string][]*WindowedElement
}

// ErrInvalidData indicates there is an issue with the Data in the GetWindowedCalculation call
var ErrInvalidData = errors.New("invalid data provided for window calculation")

// ErrMissingInformationType indicates an invalid InformationType was requested to be used to extract time series values
var ErrMissingInformationType = errors.New("missing information type in data for window calculation")

// ErrInvalidWindowLength indicates that the window length cannot be accommodated with the Data provided
var ErrInvalidWindowLength = errors.New("window length must be greater than zero and less than or equal to the length of the time series")

// WindowFunc describes the func type used by GetWindowedCalculations
type WindowFunc func(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement

// WindowedCalculationOptions provides a mechanism to alter the processing of GetWindowedCalculations
type WindowedCalculationOptions struct {
	// ElementProcessingLimit restricts GetWindowedCalculations to the specified number of Elements in the supplied Time Series
	// If not specified, then all Elements are processed
	ElementProcessingLimit int
}

var defaultWindowedCalculationOptions = WindowedCalculationOptions{
	ElementProcessingLimit: 0,
}

// WithElementProcessingLimit allows the number of Elements to be limited to the number specified
func WithElementProcessingLimit(n int) func(*WindowedCalculationOptions) error {
	return func(wco *WindowedCalculationOptions) error {
		if n < 0 {
			return fmt.Errorf("invalid number of data points: %d", n)
		}
		wco.ElementProcessingLimit = n
		return nil
	}
}

// GetWindowedCalculation performs the specified WindowFunc calculations on the supplied time series data,
// using the specified window length and informtaion type from the time series.
// Calculations are aborted if the context is ended during processing.
func GetWindowedCalculation(ctx context.Context, data *Data, windowLength int, it InformationType, calcMap map[string]WindowFunc, opts ...func(*WindowedCalculationOptions) error) (*WindowedResult, error) {

	if data == nil || !data.isValid() {
		return nil, ErrInvalidData
	}
	if !it.isValid() {
		return nil, common.ErrInvalidInformationType
	}
	if !slices.Contains(data.Meta.Information, it) {
		return nil, ErrMissingInformationType
	}
	if windowLength < 1 || windowLength > len(data.TimeSeries) {
		return nil, ErrInvalidWindowLength
	}

	result := &WindowedResult{
		Meta: &WindowedMeta{
			Input:           data,
			WindowLength:    windowLength,
			InformationType: it,
			Calculations:    calcMap,
		},
		TimeSeries: map[string][]*WindowedElement{},
	}

	if len(calcMap) == 0 {
		return result, nil
	}

	var o = defaultWindowedCalculationOptions
	o.ElementProcessingLimit = len(data.TimeSeries) // Process all data by default

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	for key, calc := range calcMap {
		if calc == nil {
			return nil, errors.New("calculation function is nil for " + key)
		}
		result.TimeSeries[key] = []*WindowedElement{}
	}

	// Make sure that the number of data points does not mean we walk off the end of the time series
	if o.ElementProcessingLimit > len(data.TimeSeries)-windowLength {
		o.ElementProcessingLimit = len(data.TimeSeries) - windowLength
	}
	result.Meta.Options = &o

	for i := range o.ElementProcessingLimit {
		select {
		case <-ctx.Done():
			return nil, common.ErrContextEnded
		default:
			for key, calc := range calcMap {
				we := calc(ctx, data, i, windowLength, it)
				result.TimeSeries[key] = append(result.TimeSeries[key], we)
			}
		}
	}

	return result, nil
}

// WindowAverage generates a time series of the mean of the value of the specified InformationType, with the
// mean calculated across the specified windowLen number of Elements at each step
func WindowAverage(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	avg := 0.0
	for j := range windowLen {
		v, ok := data.TimeSeries[offset+j].Data[it]
		if !ok {
			continue // Should not happen, but just in case
		}
		avg += v
	}

	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       avg / float64(windowLen),
	}
}

// WindowVariance generates a time series of the variance of the value of the specified InformationType, with the
// mean calculated across the specified windowLen number of Elements at each step
func WindowVariance(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	tot := 0.0
	sq := 0.0
	for j := range windowLen {
		v, ok := data.TimeSeries[offset+j].Data[it]
		if !ok {
			continue // Should not happen, but just in case
		}
		tot += v
		sq += v * v
	}
	avg := tot / float64(windowLen)
	avgSq := sq / float64(windowLen)

	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       avgSq - avg*avg,
	}
}

// WindowPercentageChange generates a time series of the percent change in value of the specified InformationType
// between the start and end of the windowLen (growth over time is positive, decline negative)
func WindowPercentageChange(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       100 * (data.TimeSeries[offset].Data[it]/data.TimeSeries[offset+windowLen].Data[it] - 1.0),
	}
}

// WindowChange generates a time series of the actual change in value of the specified InformationType
// between the start and end of the windowLen (growth over time is positive, decline negative)
func WindowChange(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       data.TimeSeries[offset].Data[it] - data.TimeSeries[offset+windowLen].Data[it],
	}
}
