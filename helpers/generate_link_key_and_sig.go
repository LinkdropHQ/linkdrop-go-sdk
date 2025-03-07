package helpers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func LinkSignatureTypedData(
	linkKeyId common.Address,
	transferId common.Address,
	domain apitypes.TypedDataDomain,
) apitypes.TypedData {
	return apitypes.TypedData{
		Domain:      domain,
		PrimaryType: "Transfer",
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Transfer": {
				{Name: "linkKeyId", Type: "address"},
				{Name: "transferId", Type: "address"},
			},
		},
		Message: map[string]any{
			"linkKeyId":  linkKeyId.Hex(),
			"transferId": transferId.Hex(),
		},
	}
}
