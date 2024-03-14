package operator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSupportStrategyExecute(t *testing.T) {
	res, err := GetSupportStrategyExecute("network")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	fmt.Println(res["1"])
}
