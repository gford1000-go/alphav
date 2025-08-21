package historic

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/gford1000-go/alphav/common"
)

type WindowedMeta struct {
	Input        *Data
	N            int
	T            InformationType
	Calculations map[string]WindowFunc
}

type WindowedElement struct {
	WindowStart time.Time
	Value       float64
}

type WindowedResult struct {
	Meta       *WindowedMeta
	TimeSeries map[string][]*WindowedElement
}

var ErrInvalidData = errors.New("invalid data provided for window calculation")

var ErrMissingInformationType = errors.New("missing information type in data for window calculation")

var ErrInvalidNEles = errors.New("neles must be greater than zero and less than or equal to the length of the time series")

type WindowFunc func(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement

type WindowedCalculationOptions struct {
	NumberOfDataPoints int
}

var defaultWindowedCalculationOptions = WindowedCalculationOptions{
	NumberOfDataPoints: 0,
}

func WithNumberOfDataPoints(n int) func(*WindowedCalculationOptions) error {
	return func(wco *WindowedCalculationOptions) error {
		if n < 0 {
			return fmt.Errorf("invalid number of data points: %d", n)
		}
		wco.NumberOfDataPoints = n
		return nil
	}
}

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
		return nil, ErrInvalidNEles
	}

	result := &WindowedResult{
		Meta: &WindowedMeta{
			Input:        data,
			N:            windowLength,
			T:            it,
			Calculations: calcMap,
		},
		TimeSeries: map[string][]*WindowedElement{},
	}

	if len(calcMap) == 0 {
		return result, nil
	}

	var o = defaultWindowedCalculationOptions
	o.NumberOfDataPoints = len(data.TimeSeries) // Process all data by default

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
	var iterations = o.NumberOfDataPoints
	if iterations > len(data.TimeSeries)-windowLength {
		iterations = len(data.TimeSeries) - windowLength
	}

	for i := 0; i < iterations; i++ {
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

func WindowAverage(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	avg := 0.0
	for j := 0; j < windowLen; j++ {
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

func WindowVariation(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	tot := 0.0
	sq := 0.0
	for j := 0; j < windowLen; j++ {
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

func WindowPercentageGrowth(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       100 * (data.TimeSeries[offset].Data[it]/data.TimeSeries[offset+windowLen].Data[it] - 1.0),
	}
}

func WindowChange(ctx context.Context, data *Data, offset, windowLen int, it InformationType) *WindowedElement {
	return &WindowedElement{
		WindowStart: data.TimeSeries[offset].Date,
		Value:       data.TimeSeries[offset].Data[it] - data.TimeSeries[offset+windowLen].Data[it],
	}
}
