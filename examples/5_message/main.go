package main

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/LinkdropHQ/linkdrop-go-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"log"
	"math/big"
	"os"
)

func signTypedData(data apitypes.TypedData) ([]byte, error) {
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}

	return utils.SignTypedData(data, privateKey)
}

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		types.DeploymentCBW,
		utils.GetRandomBytes,
		linkdrop.WithApiKey(os.Getenv("LINKDROP_API_KEY")),
		linkdrop.WithMessageConfig(
			linkdrop.MessageConfig{
				MinEncryptionKeyLength: 64,
				MaxEncryptionKeyLength: 128,
				MaxTextLength:          1000,
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	link, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeNative,
			ChainId: types.ChainIdBase,
		},
		big.NewInt(1000000000),
		common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
		big.NewInt(1000000000),
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = link.AddMessage(
		"Stay Based!",
		64,
		signTypedData,
	)
	if err != nil {
		log.Fatalln(err)
	}

	decryptedMessage, err := link.DecryptSenderMessage(signTypedData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Decrypted message: ", decryptedMessage)

	l, _, err := link.GenerateClaimUrl(signTypedData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Link with message: ", l)
}
