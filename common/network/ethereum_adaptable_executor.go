package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/contract_entrypoint_v06"
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/contract_entrypoint_v07"
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/l1_gas_oracle"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/xerrors"
	"math/big"
	"sync"
)

var PreVerificationGas = new(big.Int).SetInt64(21000)

// GweiFactor Each gwei is equal to one-billionth of an ETH (0.000000001 ETH or 10-9 ETH).
var GweiFactor = new(big.Float).SetInt(big.NewInt(1e9))
var EthWeiFactor = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
var once sync.Once
var executorMap map[types.Network]*EthereumExecutor = make(map[types.Network]*EthereumExecutor)
var TokenContractCache map[*common.Address]*contract_erc20.Contract
var V06EntryPointContractCache map[types.Network]map[common.Address]*contract_entrypoint_v06.Contract
var V07EntryPointContractCache map[types.Network]map[common.Address]*contract_entrypoint_v07.Contract

func init() {
	TokenContractCache = make(map[*common.Address]*contract_erc20.Contract)
	V06EntryPointContractCache = make(map[types.Network]map[common.Address]*contract_entrypoint_v06.Contract)
	V07EntryPointContractCache = make(map[types.Network]map[common.Address]*contract_entrypoint_v07.Contract)
}

type EthereumExecutor struct {
	BaseExecutor
	Client  *ethclient.Client
	network types.Network
}

func GetEthereumExecutor(network types.Network) *EthereumExecutor {
	once.Do(func() {
		if executorMap[network] == nil {
			client, err := ethclient.Dial(conf.GetEthereumRpcUrl(network))
			if err != nil {
				panic(err)
			}

			executorMap[network] = &EthereumExecutor{
				network: network,
				Client:  client,
			}
		}
	})
	return executorMap[network]
}

func (executor EthereumExecutor) GetUserTokenBalance(userAddress common.Address, token types.TokenType) (*big.Int, error) {
	tokenAddress := conf.GetTokenAddress(executor.network, token) //TODO
	ethTokenAddress := common.HexToAddress(tokenAddress)
	tokenInstance, err := executor.GetTokenContract(&ethTokenAddress)
	if err != nil {
		return nil, err
	}
	return tokenInstance.BalanceOf(&bind.CallOpts{}, userAddress)
}
func (executor EthereumExecutor) CheckContractAddressAccess(contract *common.Address) (bool, error) {
	client := executor.Client

	code, err := client.CodeAt(context.Background(), *contract, nil)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, xerrors.Errorf("contract  [%s] address not exist in [%s] network", contract, executor.network)
	}
	return true, nil
}

func (executor EthereumExecutor) GetTokenContract(tokenAddress *common.Address) (*contract_erc20.Contract, error) {
	client := executor.Client
	contract, ok := TokenContractCache[tokenAddress]
	if !ok {
		erc20Contract, err := contract_erc20.NewContract(*tokenAddress, client)
		if err != nil {
			return nil, err
		}
		TokenContractCache[tokenAddress] = erc20Contract
		return erc20Contract, nil
	}
	return contract, nil
}

func (executor EthereumExecutor) EstimateUserOpCallGas(entrypointAddress *common.Address, userOpParam *userop.BaseUserOp) (*big.Int, error) {
	client := executor.Client
	userOpValue := *userOpParam
	userOpValue.GetSender()
	res, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *entrypointAddress,
		To:   userOpValue.GetSender(),
		Data: userOpValue.GetCallData(),
	})
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(res), nil
}
func (executor EthereumExecutor) EstimateCreateSenderGas(entrypointAddress *common.Address, userOpParam *userop.BaseUserOp) (*big.Int, error) {
	client := executor.Client
	userOpValue := *userOpParam
	userOpValue.GetSender()
	res, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *entrypointAddress,
		To:   userOpValue.GetFactoryAddress(),
		Data: userOpValue.GetInitCode(),
	})
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(res), nil
}

func (executor EthereumExecutor) GetGasPrice() (*model.GasPrice, error) {

	client := executor.Client

	priceWei, priceWeiErr := client.SuggestGasPrice(context.Background())
	if priceWeiErr != nil {
		return nil, priceWeiErr
	}
	priorityPriceWei, tiperr := client.SuggestGasTipCap(context.Background())
	if tiperr != nil {
		return nil, tiperr
	}
	result := model.GasPrice{}
	result.MaxFeePerGas = priceWei
	result.MaxPriorityPriceWei = priorityPriceWei
	return &result, nil
	//
	//gasPriceInGwei := new(big.Float).SetInt(priceWei)
	//gasPriceInGwei.Quo(gasPriceInGwei, GweiFactor)
	//gasPriceInEther := new(big.Float).SetInt(priceWei)
	//gasPriceInEther.Quo(gasPriceInEther, EthWeiFactor)
	//gasPriceInGweiFloat, _ := gasPriceInGwei.Float64()
	//result.MaxBasePriceGwei = gasPriceInGweiFloat
	//result.MaxBasePriceEther = gasPriceInEther
	//
	//priorityPriceInGwei := new(big.Float).SetInt(priorityPriceWei)
	//priorityPriceInGwei.Quo(priorityPriceInGwei, GweiFactor)
	//priorityPriceInEther := new(big.Float).SetInt(priorityPriceWei)
	//priorityPriceInEther.Quo(priorityPriceInEther, EthWeiFactor)
	//priorityPriceInGweiFloat, _ := priorityPriceInGwei.Float64()
	//result.MaxPriorityPriceGwei = priorityPriceInGweiFloat
	//result.MaxPriorityPriceEther = gasPriceInEther
	//return &result, nil
}
func (executor EthereumExecutor) GetPreVerificationGas() (uint64, error) {
	if conf.ArbStackNetWork.Contains(executor.network) {
		return 0, nil
	}
	if conf.OpeStackNetWork.Contains(executor.network) {
		return 0, nil
	}
	return PreVerificationGas.Uint64(), nil
}

// GetL1DataFee
// OpSrource https://github.com/ethereum-optimism/optimism/blob/233ede59d16cb01bdd8e7ff662a153a4c3178bdd/packages/contracts/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L109-L124
// l1Gas = zeros * TX_DATA_ZERO_GAS + (nonzeros + 4) * TX_DATA_NON_ZERO_GAS
// l1GasFee = ((l1Gas + overhead) * l1BaseFee * scalar) / PRECISION
func (executor EthereumExecutor) GetL1DataFee(data []byte) (*big.Int, error) {
	address, ok := conf.L1GasOracleInL2[executor.network]
	if !ok {
		return nil, xerrors.Errorf("L1GasOracleInL2 not found in network %s", executor.network)
	}

	contract, err := l1_gas_oracle.NewContract(address, executor.Client)
	if err != nil {
		return nil, err
	}

	abi, err := l1_gas_oracle.ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	method := abi.Methods["getL1Fee"]
	input, err := method.Inputs.Pack(data)
	if err != nil {
		return nil, err
	}
	fee, err := contract.GetL1Fee(nil, input)
	if err != nil {
		return nil, err
	}
	return fee, nil
}

func (executor EthereumExecutor) SimulateV06HandleOp(v06 *userop.UserOperationV06, entryPoint *common.Address) (*model.SimulateHandleOpResult, error) {
	_, err := executor.GetEntryPoint06(entryPoint)
	if err != nil {
		return nil, err
	}
	//TODO
	//contract.SimulateHandleOp(nil, v06.Target, v06.Data
	//contract.SimulateHandleOp()
	return nil, nil
}
func (executor EthereumExecutor) SimulateV07HandleOp(v07 *userop.UserOperationV07, entryPoint *common.Address) (*model.SimulateHandleOpResult, error) {
	_, err := executor.GetEntryPoint07(entryPoint)
	if err != nil {
		return nil, err
	}
	//TODO
	return nil, nil
}

func (executor EthereumExecutor) GetEntryPoint07(entryPoint *common.Address) (*contract_entrypoint_v07.Contract, error) {
	contract, ok := V07EntryPointContractCache[executor.network][*entryPoint]
	if !ok {
		contractInstance, err := contract_entrypoint_v07.NewContract(*entryPoint, executor.Client)
		if err != nil {
			return nil, err
		}
		V07EntryPointContractCache[executor.network][*entryPoint] = contractInstance
		return contractInstance, nil
	}
	return contract, nil
}
func (executor EthereumExecutor) GetEntryPoint06(entryPoint *common.Address) (*contract_entrypoint_v06.Contract, error) {
	contract, ok := V06EntryPointContractCache[executor.network][*entryPoint]
	if !ok {
		contractInstance, err := contract_entrypoint_v06.NewContract(*entryPoint, executor.Client)
		if err != nil {
			return nil, err
		}
		V06EntryPointContractCache[executor.network][*entryPoint] = contractInstance
		return contractInstance, nil
	}
	return contract, nil

}
