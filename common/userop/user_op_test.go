package userop

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"testing"
)

func TestNewUserOpV06(t *testing.T) {
	userOpV6 := utils.GenerateMockUservOperation()
	userOp, err := NewUserOp(userOpV6, types.EntryPointV07)

	if err != nil {
		t.Error(err)
		return
	}
	if userOp == nil {
		t.Error("userOp is nil")
		return
	}
	//userOpvalue := *userOp
	//userOpvalueV6 := userOpvalue.(*UserOperationV07)
	//t.Logf("userOpSender: %v", userOpvalueV6.Sender)
	//t.Logf("PreVerificationGas: %v", userOpvalueV6.PreVerificationGas)

}
