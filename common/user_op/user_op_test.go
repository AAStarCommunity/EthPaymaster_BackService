package user_op

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"testing"
)

func TestNewUserOpV06(t *testing.T) {
	userOpV6 := utils.GenerateMockUservOperation()
	userOp, err := NewUserOp(userOpV6, types.EntryPointV07)
	t.Logf("userOp: %v", userOp)
	t.Logf("PreVerificationGas %v", userOp.PreVerificationGas)

	t.Logf("MaxFeePerGas %v", userOp.MaxFeePerGas)
	t.Logf("MaxPriorityFeePerGas %v", userOp.MaxPriorityFeePerGas)
	t.Logf("CallGasLimit %v", userOp.CallGasLimit)
	t.Logf("VerificationGasLimit %v", userOp.VerificationGasLimit)

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
