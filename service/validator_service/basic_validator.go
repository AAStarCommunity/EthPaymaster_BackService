package validator_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"math/big"
)

var MinPreVerificationGas *big.Int

func init() {
	MinPreVerificationGas = big.NewInt(21000)
}
func ValidateStrategy(strategy *model.Strategy) error {
	if strategy == nil {
		return xerrors.Errorf("empty strategy")
	}
	if strategy.GetNewWork() == "" {
		return xerrors.Errorf("empty strategy network")
	}
	// check Paymaster
	ok, err := chain_service.CheckContractAddressAccess(strategy.GetPaymasterAddress(), strategy.GetNewWork())
	if !ok || err != nil {
		return err
	}
	// check EntryPoint
	return nil
}

func ValidateUserOp(userOpParam *userop.BaseUserOp, strategy *model.Strategy) error {
	//recall simulate?
	//UserOp Validate
	//check nonce

	//
	if err := checkSender(userOpParam, strategy.GetNewWork()); err != nil {
		return err
	}
	userOpValue := *userOpParam
	if !userOpValue.GetNonce().IsInt64() {
		return xerrors.Errorf("nonce is not in uint64 range")
	}

	return userOpValue.ValidateUserOp()

	//If initCode is not empty, parse its first 20 bytes as a factory address. Record whether the factory is staked, in case the later simulation indicates that it needs to be. If the factory accesses global state, it must be staked - see reputation, throttling and banning section for details.
	//The verificationGasLimit is sufficiently low (<= MAX_VERIFICATION_GAS) and the preVerificationGas is sufficiently high (enough to pay for the calldata gas cost of serializing the UserOperationV06 plus PRE_VERIFICATION_OVERHEAD_GAS)
	//

	//validate trusted entrypoint
}
func checkSender(userOpParam *userop.BaseUserOp, netWork types.Network) error {
	//check sender
	userOpValue := *userOpParam
	checkOk, checkSenderErr := chain_service.CheckContractAddressAccess(userOpValue.GetSender(), netWork)
	if !checkOk {
		if err := checkInitCode(userOpValue.GetInitCode(), netWork); err != nil {
			return xerrors.Errorf("%s and %s", checkSenderErr.Error(), err.Error())
		}
	}
	//check balance
	return nil
}
func checkInitCode(initCode []byte, network types.Network) error {
	if len(initCode) < 20 {
		return xerrors.Errorf("initCode length is less than 20 do not have factory address")
	}
	factoryAddress := common.BytesToAddress(initCode[:20])
	if ok, err := chain_service.CheckContractAddressAccess(&factoryAddress, network); err != nil {
		return err
	} else if !ok {
		return xerrors.Errorf("factoryAddress address [factoryAddress] not exist in [%s] network", network)
	}
	//parse its first 20 bytes as a factory address. Record whether the factory is staked,
	//factory and factoryData - either both exist, or none

	//parse its first 20 bytes as a factory address. Record whether the factory is staked,
	//factory and factoryData - either both exist, or none
	return nil
}
