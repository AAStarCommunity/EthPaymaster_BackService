package validator_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
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

func ValidateUserOp(userOp *userop.BaseUserOp) error {
	//recall simulate?
	//UserOp Validate
	//check nonce
	//if userOp.PreVerificationGas.Cmp(MinPreVerificationGas) == -1 {
	//	return xerrors.Errorf("preVerificationGas is less than 21000")
	//}
	//
	//if err := checkSender(userOp, network.Sepolia); err != nil {
	//	return err
	//}
	//
	//if !userOp.Nonce.IsInt64() {
	//	return xerrors.Errorf("nonce is not in uint64 range")
	//}

	//If initCode is not empty, parse its first 20 bytes as a factory address. Record whether the factory is staked, in case the later simulation indicates that it needs to be. If the factory accesses global state, it must be staked - see reputation, throttling and banning section for details.
	//The verificationGasLimit is sufficiently low (<= MAX_VERIFICATION_GAS) and the preVerificationGas is sufficiently high (enough to pay for the calldata gas cost of serializing the UserOperation plus PRE_VERIFICATION_OVERHEAD_GAS)
	// check Sender is valid ,if sender is invalid And InitCode empty, return error
	//
	// nonce is valid
	//validate trusted entrypoint
	return nil
}
func checkSender(userOp *userop.UserOperation, netWork network.Network) error {
	//check sender
	checkOk, checkSenderErr := chain_service.CheckContractAddressAccess(userOp.Sender, netWork)
	if !checkOk {
		if err := checkInitCode(userOp.InitCode, netWork); err != nil {
			return xerrors.Errorf("%s and %s", checkSenderErr.Error(), err.Error())
		}
	}
	//check balance

	return nil
}
func checkInitCode(initCode []byte, network network.Network) error {
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
