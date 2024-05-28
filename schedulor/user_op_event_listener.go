package schedulor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_contract/contract/contract_entrypoint_v06"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/price_compoent"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
)

var (
	userOpEvent   abi.Event
	entryPointABI *abi.ABI
)

func init() {
	getAbi, _ := contract_entrypoint_v06.ContractMetaData.GetAbi()
	entryPointABI = getAbi
	userOpEvent = entryPointABI.Events["UserOperationEvent"]
}

type EventListener struct {
	logCh               chan types.Log
	entryPointAddresses []common.Address
	paymasterAddresses  []common.Address
	client              *ethclient.Client
	network             global_const.Network
}

func (listener *EventListener) Listen() {
	queryies := make([][]interface{}, 0)
	eventId := userOpEvent.ID
	var eventIds []interface{}
	eventIds = append(eventIds, eventId)
	queryies = append(queryies, eventIds)
	var userOpHashRule []interface{}
	queryies = append(queryies, userOpHashRule)
	var senderRule []interface{}
	queryies = append(queryies, senderRule)
	var paymasterRule []interface{}
	for _, paymasterAddress := range listener.paymasterAddresses {
		paymasterRule = append(paymasterRule, paymasterAddress)
	}
	queryies = append(queryies, paymasterRule)
	topics, err := abi.MakeTopics(queryies...)
	if err != nil {
		logrus.Errorf("abi.MakeTopics failed: %v", err)
		return
	}
	addresses := make([]common.Address, 0)
	for _, entryPointAdd := range listener.entryPointAddresses {
		addresses = append(addresses, entryPointAdd)
	}

	query := ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
	fmt.Println("query: ", query)
	client := listener.client

	sub, err := client.SubscribeFilterLogs(context.Background(), query, listener.logCh)

	if err != nil {
		logrus.Errorf("SubscribeFilterLogs failed: %v", err)
		return
	}
	for {
		select {
		case err := <-sub.Err():
			logrus.Errorf("SubscribeFilterLogs failed: %v", err)
			fmt.Println("SubscribeFilterLogs failed: ", err)
			return
		case vLog := <-listener.logCh:
			if len(vLog.Topics) == 0 {
				logrus.Errorf("vLog.Topics is empty")
				continue
			}
			if vLog.Topics[0] != eventId {
				logrus.Errorf("vLog.Topics[0] != eventId")
				continue
			}
			if len(vLog.Topics) < 4 {
				logrus.Errorf("vLog.Topics length < 4")
				continue
			}
			if len(vLog.Data) > 0 {
				eventObj := &ContractUserOperationEvent{}
				err := entryPointABI.UnpackIntoInterface(eventObj, "UserOperationEvent", vLog.Data)
				if err != nil {
					logrus.Errorf("UnpackIntoInterface failed: %v", err)
					continue
				}
				eventObj.UserOpHash = common.BytesToHash(vLog.Topics[1].Bytes())
				eventObj.Sender = common.BytesToAddress(vLog.Topics[2].Bytes())
				eventObj.Paymaster = common.BytesToAddress(vLog.Topics[3].Bytes())
				logrus.Debugf("UserOpEventComunicate: %v, %v, %v, %v", eventObj.UserOpHash, eventObj.Sender, eventObj.ActualGasCost, eventObj.ActualGasUsed)
				UserOpEventComunicate(listener.network, *eventObj)
			}
		}
	}

}

func NewEventListener(client *ethclient.Client, network global_const.Network) (EventListener, error) {
	entryPointAddresses, err := config.GetSupportEntryPoints(network)
	if err != nil {
		return EventListener{}, err
	}
	entryPointAddressesArr := make([]common.Address, 2)
	for _, address := range entryPointAddresses.ToSlice() {
		entryPointAddressesArr = append(entryPointAddressesArr, common.HexToAddress(address))
	}

	paymasterAddresses, err := config.GetSupportPaymaster(network)
	if err != nil {
		return EventListener{}, err
	}
	paymasterAddressArr := make([]common.Address, 2)
	for _, address := range paymasterAddresses.ToSlice() {
		paymasterAddressArr = append(paymasterAddressArr, common.HexToAddress(address))
	}

	return EventListener{
		logCh:               make(chan types.Log, 10),
		entryPointAddresses: entryPointAddressesArr,
		paymasterAddresses:  paymasterAddressArr,
		client:              client,
		network:             network,
	}, nil
	//if entryPointAddress == nil {
	//	logrus.Debugf("Not Support Network %v", network)
	//	return
	//}
	//for _, address := range entryPointAddress.ToSlice() {
	//	contractAddress := common.HexToAddress(address)
	//	query := ethereum.FilterQuery{
	//		Addresses: []common.Address{contractAddress},
	//	}
	//	sub, err := client.SubscribeFilterLogs(nil, query, nil)
	//	if err != nil {
	//		logrus.Errorf("SubscribeFilterLogs failed: %v", err)
	//		return
	//	}
	//	logs := make(chan types.Log)
	//	go func() {
	//		for {
	//			select {
	//			case err := <-sub.Err():
	//				logrus.Errorf("SubscribeFilterLogs failed: %v", err)
	//				return
	//			case vLog := <-logs:
	//				if vLog.Removed {
	//					continue
	//				}
	//				//TODO
	//				//UserOpEventComunicate(network, ContractUserOperationEvent{})
	//			}
	//		}
	//	}()
	//	//TODO
	//	//client.SubscribeFilterLogs()
	//}
	////TODO

}

type ContractUserOperationEvent struct {
	UserOpHash    [32]byte
	Sender        common.Address
	Paymaster     common.Address
	Nonce         *big.Int
	Success       bool
	ActualGasCost *big.Int
	ActualGasUsed *big.Int
}

func UserOpEventComunicate(network global_const.Network, event ContractUserOperationEvent) {
	paymasterAddressSet, _ := config.GetSupportPaymaster(network)
	if paymasterAddressSet == nil {
		logrus.Debugf("Not Support Network %v", network)
		return
	}
	//if !paymasterAddressSet.Contains(event.Paymaster) {
	//	logrus.Debugf("UserOpEventComunicate: paymaster not support, %v", event.Paymaster)
	//	return
	//}
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

	_, err = sponsor_manager.ReleaseBalanceWithActualCost(event.Sender.String(), event.UserOpHash[:], gasCostUsd, true)
	if err != nil {
		//TODO if is NetWorkError, need retry
		logrus.Errorf("ReleaseBalanceWithActualCost failed: %v", err)
		return
	}
}
