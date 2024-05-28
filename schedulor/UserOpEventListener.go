package schedulor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/price_compoent"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
)

func AddListener(client ethclient.Client, network global_const.Network) {
	entryPointAddress, _ := config.GetSupportEntryPoints(network)
	if entryPointAddress == nil {
		logrus.Debugf("Not Support Network %v", network)
		return
	}
	for _, address := range entryPointAddress.ToSlice() {
		contractAddress := common.HexToAddress(address)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{contractAddress},
		}
		sub, err := client.SubscribeFilterLogs(nil, query, nil)
		if err != nil {
			logrus.Errorf("SubscribeFilterLogs failed: %v", err)
			return
		}
		logs := make(chan types.Log)
		go func() {
			for {
				select {
				case err := <-sub.Err():
					logrus.Errorf("SubscribeFilterLogs failed: %v", err)
					return
				case vLog := <-logs:
					if vLog.Removed {
						continue
					}
					//TODO
					//UserOpEventComunicate(network, ContractUserOperationEvent{})
				}
			}
		}()
		//TODO
		//client.SubscribeFilterLogs()
	}
	//TODO
}

type ContractUserOperationEvent struct {
	UserOpHash    [32]byte
	Sender        string
	Paymaster     string
	Nonce         *big.Int
	Success       bool
	ActualGasCost *big.Int
	ActualGasUsed *big.Int
	Raw           types.Log
}

func UserOpEventComunicate(network global_const.Network, event ContractUserOperationEvent) {
	paymasterAddressSet, _ := config.GetSupportPaymaster(network)
	if paymasterAddressSet == nil {
		logrus.Debugf("Not Support Network %v", network)
		return
	}
	if !paymasterAddressSet.Contains(event.Paymaster) {
		logrus.Debugf("UserOpEventComunicate: paymaster not support, %v", event.Paymaster)
		return
	}
	if !event.Success {
		_, err := sponsor_manager.ReleaseUserOpHashLockWhenFail(event.UserOpHash[:], true)
		if err != nil {
			logrus.Errorf("ReleaseUserOpHashLockWhenFail failed: %v", err)
		}
		return
	}
	gasCostEther := new(big.Float).SetInt(event.ActualGasCost)
	gasCostEther = new(big.Float).Quo(gasCostEther, new(big.Float).SetInt64(1e18))
	//logrus.Infof("UserOpEventComunicate: %v, %v, %v, %v", event.UserOpHash, event.Sender, event.ActualGasCost, gasCostEther)
	gasCostUsd, err := price_compoent.GetTokenCostInUsd(global_const.TokenTypeETH, gasCostEther)
	if err != nil {
		//TODO if is NetWorkError, need retry
		logrus.Errorf("GetTokenCostInUsd failed: %v", err)
		return
	}

	err = sponsor_manager.ReleaseBalanceWithActualCost(event.Sender, event.UserOpHash[:], network, gasCostUsd, true)
	if err != nil {
		//TODO if is NetWorkError, need retry
		logrus.Errorf("ReleaseBalanceWithActualCost failed: %v", err)
		return
	}
}
