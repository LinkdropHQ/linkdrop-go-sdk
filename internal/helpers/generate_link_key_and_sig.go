package helpers

import (
	"crypto/ecdsa"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func GenerateLinkKey(
	getRandomBytes types.RandomBytesCallback,
) (linkKey *ecdsa.PrivateKey, linkKeyId common.Address, err error) {
	linkKey, err = PrivateKey(getRandomBytes)
	if err != nil {
		return
	}
	linkKeyId = crypto.PubkeyToAddress(linkKey.PublicKey)
	return
}

// GenerateLinkSignature
// NOTE: Usually will be called after GenerateLinkKey if signature is needed
func GenerateLinkSignature(
	linkKeyId common.Address,
	signTypedData types.SignTypedDataCallback,
	transferId common.Address,
	domain apitypes.TypedDataDomain,
) (senderSig []byte, err error) {
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
			"linkKeyId":  linkKeyId.Hex(),
			"transferId": transferId.Hex(),
		},
	})
	return
}
