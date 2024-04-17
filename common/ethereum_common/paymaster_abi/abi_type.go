package paymaster_abi

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	BytesType   abi.Type
	Bytes32Type abi.Type
	Uint48Type  abi.Type
	Uint256Type abi.Type
	Uint64Type  abi.Type
	AddressType abi.Type
	BooleanType abi.Type
)

func init() {
	bytesTypeVar, err := abi.NewType("bytes", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	BytesType = bytesTypeVar

	byte32TypeVar, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	Bytes32Type = byte32TypeVar

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

	uint64TypeVar, err := abi.NewType("uint64", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	Uint64Type = uint64TypeVar

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	AddressType = addressType

	boolTypeVar, err := abi.NewType("bool", "", nil)
	if err != nil {
		panic(fmt.Sprintf("[initerror] %s", err))
	}
	BooleanType = boolTypeVar
}
