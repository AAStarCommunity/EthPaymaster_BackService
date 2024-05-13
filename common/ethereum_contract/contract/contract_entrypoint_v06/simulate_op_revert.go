package contract_entrypoint_v06

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_contract/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math/big"
)

type ExecutionResultRevert struct {
	PreOpGas      *big.Int
	Paid          *big.Int
	ValidAfter    *big.Int
	ValidUntil    *big.Int
	TargetSuccess bool
	TargetResult  []byte
}

func ExecutionResult() abi.Error {
	return abi.NewError("ExecutionResult", abi.Arguments{
		{Name: "preOpGas", Type: paymaster_abi.Uint256Type},
		{Name: "paid", Type: paymaster_abi.Uint256Type},
		{Name: "validAfter", Type: paymaster_abi.Uint48Type},
		{Name: "validUntil", Type: paymaster_abi.Uint48Type},
		{Name: "targetSuccess", Type: paymaster_abi.BooleanType},
		{Name: "targetResult", Type: paymaster_abi.BytesType},
	})
}

func NewExecutionResult(inputError error, abi *abi.ABI) (*ExecutionResultRevert, error) {

	var rpcErr rpc.DataError
	ok := errors.As(inputError, &rpcErr)
	if !ok {
		return nil, xerrors.Errorf("ExecutionResult: cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)

	logrus.Debug("data: ", data)
	if !ok {
		return nil, xerrors.Errorf("ExecutionResult: cannot assert type: data is not of type string")
	}

	sim := ExecutionResult()
	revert, upPackerr := sim.Unpack(common.Hex2Bytes(data[2:]))
	if upPackerr != nil {
		// is another ERROR
		errstr, parseErr := utils.ParseCallError(inputError, abi)
		return nil, fmt.Errorf("executionResult err: [%s] parErr [%s]", errstr, parseErr)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, xerrors.Errorf("executionResult: cannot assert type: args is not of type []any")
	}
	if len(args) != 6 {
		return nil, xerrors.Errorf("executionResult: invalid args length: expected 6, got %d", len(args))
	}

	return &ExecutionResultRevert{
		PreOpGas:      args[0].(*big.Int),
		Paid:          args[1].(*big.Int),
		ValidAfter:    args[2].(*big.Int),
		ValidUntil:    args[3].(*big.Int),
		TargetSuccess: args[4].(bool),
		TargetResult:  args[5].([]byte),
	}, nil
}
