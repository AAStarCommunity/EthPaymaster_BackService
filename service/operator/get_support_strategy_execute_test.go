package operator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSupportStrategyExecute(t *testing.T) {
	res, err := GetSupportStrategyExecute("network")
	assert.NoError(t, err)
	assert.NotNil(t, res)

}
