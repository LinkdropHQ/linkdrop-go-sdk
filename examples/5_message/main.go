package main

import (
	"crypto/rand"
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

func getRandomBytes(length int64) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return b
}

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
		os.Getenv("LINKDROP_API_KEY"),
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

	link, err := sdk.ClaimLink(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdBase,
				Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
			},
			Sender:     common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
			Amount:     big.NewInt(100000),
			Expiration: 1773159165,
		},
		getRandomBytes,
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = link.AddMessage(
		"Stay Based!",
		64,
		signTypedData,
		getRandomBytes,
	)
	if err != nil {
		log.Fatalln(err)
	}

	decryptedMessage, err := link.DecryptSenderMessage(signTypedData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Decrypted message: ", decryptedMessage)

	l, err := link.GenerateClaimUrl(nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Link with message: ", l)
}
