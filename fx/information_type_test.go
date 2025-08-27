package fx

import "testing"

func TestInformationType(t *testing.T) {

	type test struct {
		v           InformationType
		shouldPanic bool
	}

	tests := []test{
		{
			v: Open,
		},
		{
			v: High,
		},
		{
			v: Low,
		},
		{
			v: Close,
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

	runTest := func(v InformationType, shouldPanic bool) {
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
