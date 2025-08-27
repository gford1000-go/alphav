package fx

import "testing"

func TestIntradayInformationType(t *testing.T) {

	type test struct {
		v           IntradayInformationType
		shouldPanic bool
	}

	tests := []test{
		{
			v: FXRate,
		},
		{
			v: Bid,
		},
		{
			v: Ask,
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

	runTest := func(v IntradayInformationType, shouldPanic bool) {
		var panicked = new(bool)

		test := func() {
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

			_ = v.String()
		}

		test()
	}

	for _, tst := range tests {
		runTest(tst.v, tst.shouldPanic)
	}
}
