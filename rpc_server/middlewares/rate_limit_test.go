package middlewares

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestRateLimitShouldPreventRequestWhenOverDefaultLimit",
			func(t *testing.T) {
				testRateLimitShouldPreventRequestWhenOverDefaultLimit(t)
			},
		},
		{
			"TestRateLimiterShouldAllowDefaultLimitPerSecond",
			func(t *testing.T) {
				testRateLimiterShouldAllowDefaultLimitPerSecond(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func testRateLimitShouldPreventRequestWhenOverDefaultLimit(t *testing.T) {

	mockApiKey := "TestingAipKey"

	// assuming this for loop taking less than 1 second to finish
	for i := 0; i < int(DefaultLimit)+5; i++ {
		b := limiting(&mockApiKey)
		if i < int(DefaultLimit) {
			assert.Equal(t, true, b)
		} else {
			assert.Equal(t, false, b)
		}
	}
	clearLimiter(&mockApiKey)
}

func testRateLimiterShouldAllowDefaultLimitPerSecond(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Logf("Skip test in GitHub Actions")
		return
	}
	mockApiKey := "TestingAipKey"

	for x := 1; x <= 2; x++ {
		for i := 0; i < int(DefaultLimit)+5; i++ {
			b := limiting(&mockApiKey)
			if i < int(DefaultLimit) {
				assert.Equal(t, true, b)
			} else {
				assert.Equal(t, false, b)
			}
		}
		time.Sleep(time.Second * 2)
	}
	clearLimiter(&mockApiKey)
}
