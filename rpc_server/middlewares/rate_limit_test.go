package middlewares

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRateLimitShouldPreventRequestWhenOverDefaultLimit(t *testing.T) {

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
}

func TestRateLimiterShouldAllowDefaultLimitPerSecond(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip()
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
}
