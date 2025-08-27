package common

import "testing"

func TestRoundTo(t *testing.T) {

	type test struct {
		value  float64
		result float64
		dp     int
	}

	tests := []test{
		{
			value:  1,
			result: 1,
			dp:     2,
		},
		{
			value:  1.23,
			result: 1.23,
			dp:     2,
		},
		{
			value:  1.23456,
			result: 1.23,
			dp:     2,
		},
		{
			value:  1.23456,
			result: 1.235,
			dp:     3,
		},
		{
			value:  1.23456,
			result: 1.2346,
			dp:     4,
		},
		{
			value:  1.23456,
			result: 1.23456,
			dp:     5,
		},
		{
			value:  1.230456,
			result: 1.23,
			dp:     3,
		},
	}

	for i, tt := range tests {
		r := RoundTo(tt.value, tt.dp)
		if r != tt.result {
			t.Fatalf("%d: mismatch for %v", i, r)
		}
	}
}
