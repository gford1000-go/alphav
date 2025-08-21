package alphav

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// *** Add your Alpha Vantage TEST KEY here ***
	// Note that free api keys are limited to 25 requests/day.  This is based on IP address not the key itself.
	// Higher use requires premium access: see https://www.alphavantage.co/premium/
	os.Setenv("AV_API_KEY", "ATESTAPIKEY")

	// Run tests
	code := m.Run()

	// Clean up if needed
	os.Unsetenv("AV_API_KEY")

	os.Exit(code)
}
