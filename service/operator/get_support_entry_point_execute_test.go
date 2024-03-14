package operator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSupportEntrypointExecute(t *testing.T) {
	res, err := GetSupportEntrypointExecute("network")
	assert.NoError(t, err)
	t.Log(res)
}
