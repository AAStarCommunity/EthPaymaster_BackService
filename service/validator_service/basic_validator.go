package validator_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
)

func ValidateStrategy(strategy *model.Strategy, userOp *model.UserOperationItem) error {
	if strategy == nil {
		return xerrors.Errorf("empty strategy")
	}
	if strategy.NetWork == "" {
		return xerrors.Errorf("empty strategy network")
	}
	// check Paymaster
	ok, err := chain_service.CheckContractAddressAccess(strategy.PayMasterAddress, strategy.NetWork)
	if !ok || err != nil {
		return err
	}
	// check EntryPoint
	return nil
}

func ValidateUserOp(userOp *model.UserOperationItem) error {

	if err := checkSender(userOp, types.Sepolia); err != nil {
		return err
	}
	if !utils.IsStringInUint64Range(userOp.Nonce) {
		return xerrors.Errorf("nonce is not in uint64 range")
	}

	//If initCode is not empty, parse its first 20 bytes as a factory address. Record whether the factory is staked, in case the later simulation indicates that it needs to be. If the factory accesses global state, it must be staked - see reputation, throttling and banning section for details.
	//The verificationGasLimit is sufficiently low (<= MAX_VERIFICATION_GAS) and the preVerificationGas is sufficiently high (enough to pay for the calldata gas cost of serializing the UserOperation plus PRE_VERIFICATION_OVERHEAD_GAS)
	// check Sender is valid ,if sender is invalid And InitCode empty, return error
	//
	// nonce is valid
	//validate trusted entrypoint
	return nil
}
func checkSender(userOp *model.UserOperationItem, netWork types.Network) error {
	//check sender
	if userOp.Sender != "" {
		if ok, err := chain_service.CheckContractAddressAccess(userOp.Sender, netWork); err != nil {
			return err
		} else if !ok {
			return xerrors.Errorf("sender address not exist in [%s] network", netWork)
		}
		//check balance
	}
	if userOp.InitCode == "" {
		return xerrors.Errorf("initCode can not be empty if sender is empty")
	}
	if err := checkInitCode(userOp.InitCode, netWork); err != nil {

	}
	return nil
}
func checkInitCode(initCode string, network types.Network) error {
	initCodeByte, err := hex.DecodeString(initCode)
	if err != nil {
		return xerrors.Errorf("initCode is not hex string %s", initCode)
	}
	if len(initCodeByte) < 20 {
		return xerrors.Errorf("initCode is not valid %s", initCode)
	}
	factoryAddress := common.BytesToAddress(initCodeByte[:20])
	if ok, err := chain_service.CheckContractAddressAccess(factoryAddress.String(), network); err != nil {
		return err
	} else if !ok {
		return xerrors.Errorf("sender address not exist in [%s] network", network)
	}
	// TODO checkFactoryAddress stack
	//parse its first 20 bytes as a factory address. Record whether the factory is staked,
	//factory and factoryData - either both exist, or none

	//parse its first 20 bytes as a factory address. Record whether the factory is staked,
	//factory and factoryData - either both exist, or none
	return nil
}
