package helpers

import (
	"crypto/ecdsa"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateLinkKeyAndSignature(
	signTypedData types.SignTypedDataCallback,
	getRandomBytes types.RandomBytesCallback,
	transferId common.Address,
	domain types.TypedDataDomain,
) (linkKey *ecdsa.PrivateKey, linkKeyId common.Address, senderSig string, err error) {
	linkKey, err = PrivateKey(getRandomBytes)
	if err != nil {
		return
	}

	linkKeyId = crypto.PubkeyToAddress(linkKey.PublicKey)

	typedData := map[string][]types.TypedDataField{
		"Transfer": {
			{Name: "linkKeyId", Type: "address"},
			{Name: "transferId", Type: "address"},
		},
	}

	message := map[string]interface{}{
		"linkKeyId":  linkKeyId.Hex(),
		"transferId": transferId.Hex(),
	}

	senderSig, err = signTypedData(domain, typedData, message)

	return
}
