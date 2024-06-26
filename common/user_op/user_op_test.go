package user_op

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"testing"
)

func TestUserOp(t *testing.T) {
	userOpMap := utils.GenerateMockUservOperation()
	userOp, err := NewUserOp(userOpMap)
	if err != nil {
		t.Error(err)
		return
	}
	if userOp == nil {
		t.Error("userOp is nil")
		return
	}
	t.Logf("userOp: %v", userOp)

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestPackUserOpV6",
			func(t *testing.T) {
				testPackUserOp(t, userOp, global_const.EntrypointV06)
			},
		},
		{
			"TestPackUserOpV7",
			func(t *testing.T) {
				testPackUserOp(t, userOp, global_const.EntrypointV07)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testPackUserOp(t *testing.T, userOp *UserOpInput, version global_const.EntrypointVersion) {
	res, _, err := userOp.PackUserOpForMock(version)
	if err != nil {
		t.Error(err)
		return
	}
	if res == "" {
		t.Error("res is nil")
		return
	}

}
