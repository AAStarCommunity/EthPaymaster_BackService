generate-verifyingPaymaster-pkg:
	abigen --abi=./common/abi/verifying_paymaster_abi.json --pkg=contract --out=./common/contract/ethereum_compatible_contract
