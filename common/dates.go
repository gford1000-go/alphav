package common

import "time"

// ParseIntradayDate parses a string in the format "2006-01-02 15:04:05" into a time.Time object.
func ParseIntradayDate(tmStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", tmStr)
}

// ParseDate parses a string in the format "2006-01-02" into a time.Time object.
func ParseDate(dtStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dtStr)
}
