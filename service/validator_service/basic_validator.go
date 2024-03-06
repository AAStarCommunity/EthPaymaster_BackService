package validator_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
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
	// check Sender is valid ,if sender is invalid And InitCode empty, return error
	// nonce is valid
	return nil
}
