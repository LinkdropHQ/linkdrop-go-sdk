package helpers

import (
	"crypto/ecdsa"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func GenerateLinkKeyAndSignature(
	signTypedData types.SignTypedDataCallback,
	getRandomBytes types.RandomBytesCallback,
	transferId common.Address,
	domain apitypes.TypedDataDomain,
) (linkKey *ecdsa.PrivateKey, linkKeyId common.Address, senderSig []byte, err error) {

	linkKey, err = PrivateKey(getRandomBytes)
	if err != nil {
		return
	}

	senderSig, err = signTypedData(apitypes.TypedData{
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
			"linkKeyId":  crypto.PubkeyToAddress(linkKey.PublicKey).Hex(),
			"transferId": transferId.Hex(),
		},
	})
	return
}
