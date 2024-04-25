package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var GenerateFuncMap = map[global_const.PayType]GeneratePaymasterDataFunc{}
var BasicPaymasterDataAbi abi.Arguments

func init() {
	GenerateFuncMap[global_const.PayTypeVerifying] = GenerateBasicPaymasterData()
	GenerateFuncMap[global_const.PayTypeERC20] = GenerateBasicPaymasterData()
	GenerateFuncMap[global_const.PayTypeSuperVerifying] = GenerateSuperContractPaymasterData()
	BasicPaymasterDataAbi = getAbiArgs()
}
func GetGenerateFunc(payType global_const.PayType) GeneratePaymasterDataFunc {
	return GenerateFuncMap[payType]
}
func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		{Name: "validUntil", Type: paymaster_abi.Uint48Type},
		{Name: "validAfter", Type: paymaster_abi.Uint48Type},
		{Name: "erc20Token", Type: paymaster_abi.AddressType},
		{Name: "exchangeRate", Type: paymaster_abi.Uint256Type},
	}
}

type GeneratePaymasterDataFunc = func(data *paymaster_data.PaymasterData, signature []byte) ([]byte, error)

func GenerateBasicPaymasterData() GeneratePaymasterDataFunc {
	return func(data *paymaster_data.PaymasterData, signature []byte) ([]byte, error) {
		packed, err := BasicPaymasterDataAbi.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
		if err != nil {
			return nil, err
		}
		concat := data.Paymaster.Bytes()
		concat = append(concat, packed...)
		concat = append(concat, signature...)
		return concat, nil
	}
}

func GenerateSuperContractPaymasterData() GeneratePaymasterDataFunc {
	return func(data *paymaster_data.PaymasterData, signature []byte) ([]byte, error) {
		packed, err := BasicPaymasterDataAbi.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
		if err != nil {
			return nil, err
		}

		concat := data.Paymaster.Bytes()
		concat = append(concat, packed...)
		concat = append(concat, signature...)
		return concat, nil
	}
}
