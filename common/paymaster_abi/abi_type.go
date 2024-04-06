package paymaster_abi

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	BytesType   abi.Type
	Uint48Type  abi.Type
	Uint256Type abi.Type
	AddressType abi.Type
)

func init() {
	bytesTypeVar, err := abi.NewType("bytes", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	BytesType = bytesTypeVar

	uint48Var, err := abi.NewType("uint48", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	Uint48Type = uint48Var

	uint256TypeVar, err := abi.NewType("uint256", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	Uint256Type = uint256TypeVar

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	AddressType = addressType
}
