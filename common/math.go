package common

import "math"

// RoundTo will round the specified float to the specified number of decimal places
func RoundTo(x float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(x*pow) / pow
}

// EqualFloat64 returns true if the provided float64 are equal to the specified number of decimal places
func EqualFloat64(a, b float64, n int) bool {
	return RoundTo(a, n) == RoundTo(b, n)
}
