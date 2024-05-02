package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"golang.org/x/xerrors"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var GenerateFuncMap = map[global_const.PayType]GeneratePaymasterDataFunc{}
var BasicPaymasterDataAbiV06 abi.Arguments
var BasicPaymasterDataAbiV07 abi.Arguments

func init() {
	GenerateFuncMap[global_const.PayTypeVerifying] = GenerateBasicPaymasterData()
	GenerateFuncMap[global_const.PayTypeERC20] = GenerateBasicPaymasterData()
	GenerateFuncMap[global_const.PayTypeSuperVerifying] = GenerateSuperContractPaymasterData()
	BasicPaymasterDataAbiV07 = abi.Arguments{
		{Name: "accountGasLimit", Type: paymaster_abi.Bytes32Type},
		{Name: "validUntil", Type: paymaster_abi.Uint48Type},
		{Name: "validAfter", Type: paymaster_abi.Uint48Type},
		{Name: "erc20Token", Type: paymaster_abi.AddressType},
		{Name: "exchangeRate", Type: paymaster_abi.Uint256Type},
	}
	BasicPaymasterDataAbiV06 = abi.Arguments{
		{Name: "validUntil", Type: paymaster_abi.Uint48Type},
		{Name: "validAfter", Type: paymaster_abi.Uint48Type},
		{Name: "erc20Token", Type: paymaster_abi.AddressType},
		{Name: "exchangeRate", Type: paymaster_abi.Uint256Type},
	}

}
func GetGenerateFunc(payType global_const.PayType) GeneratePaymasterDataFunc {
	return GenerateFuncMap[payType]
}

type GeneratePaymasterDataFunc = func(data *paymaster_data.PaymasterDataInput, signature []byte) ([]byte, error)

func GenerateBasicPaymasterData() GeneratePaymasterDataFunc {
	return func(data *paymaster_data.PaymasterDataInput, signature []byte) ([]byte, error) {
		var packedRes []byte
		if data.EntryPointVersion == global_const.EntrypointV06 {
			v06Packed, err := BasicPaymasterDataAbiV06.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
			if err != nil {
				return nil, err
			}
			packedRes = v06Packed
		} else if data.EntryPointVersion == global_const.EntrypointV07 {
			accountGasLimit := utils.PackIntTo32Bytes(data.PaymasterVerificationGasLimit, data.PaymasterPostOpGasLimit)
			v07Packed, err := BasicPaymasterDataAbiV07.Pack(accountGasLimit, data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
			if err != nil {
				return nil, err
			}
			packedRes = v07Packed
		} else {
			return nil, xerrors.Errorf("unsupported entrypoint version")
		}

		concat := data.Paymaster.Bytes()
		concat = append(concat, packedRes...)
		concat = append(concat, signature...)
		return concat, nil
	}
}

func GenerateSuperContractPaymasterData() GeneratePaymasterDataFunc {
	return func(data *paymaster_data.PaymasterDataInput, signature []byte) ([]byte, error) {
		packed, err := BasicPaymasterDataAbiV06.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
		if err != nil {
			return nil, err
		}

		concat := data.Paymaster.Bytes()
		concat = append(concat, packed...)
		concat = append(concat, signature...)
		return concat, nil
	}
}
