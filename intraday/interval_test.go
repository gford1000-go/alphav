package intraday

import "testing"

func TestInterval(t *testing.T) {

	type test struct {
		v           Interval
		shouldPanic bool
	}

	tests := []test{
		{
			v: OneMin,
		},
		{
			v: FiveMin,
		},
		{
			v: FifteenMin,
		},
		{
			v: ThirtyMin,
		},
		{
			v: SixtyMin,
		},
		{
			v:           0,
			shouldPanic: true,
		},
		{
			v:           -99,
			shouldPanic: true,
		},
		{
			v:           99,
			shouldPanic: true,
		},
	}

	runTest := func(v Interval, shouldPanic bool) {
		var panicked = new(bool)

		test := func() string {
			defer func() {
				if (shouldPanic && !*panicked) || (!shouldPanic && *panicked) {
					t.Fatalf("unexpected error for %d", v)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					*panicked = true
				}
			}()

			return v.String()
		}

		s := test()
		if s == "" {
			return // This occurs for panicking tests
		}

		v1, err := parseInterval(s)
		if err != nil {
			t.Fatalf("unexpected parse failure for %s", s)
		}
		if v != v1 {
			t.Fatalf("unexpected parse output: expected: %v, got: %v", v, v1)
		}
	}

	for _, tst := range tests {
		runTest(tst.v, tst.shouldPanic)
	}
}
