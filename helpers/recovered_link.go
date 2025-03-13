package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func RecoveredLinkTypedData(
	linkKeyId common.Address,
	transferId common.Address,
	chainId types.ChainId,
	escrowVersion string,
	escrow common.Address,
) apitypes.TypedData {
	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:              "LinkdropEscrow",
			Version:           escrowVersion,
			ChainId:           math.NewHexOrDecimal256(int64(chainId)),
			VerifyingContract: escrow.Hex(),
		},
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
