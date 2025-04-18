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
	"os"
)

func signTypedData(data apitypes.TypedData) ([]byte, error) {
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}

	return utils.SignTypedData(data, privateKey)
}

func getRandomBytes(length int64) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return b
}

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
		linkdrop.WithApiUrl("https://escrow-api.linkdrop.io/v3"),
		linkdrop.WithProductionDefaults(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	clRecovered, err := sdk.ClaimLinkRecovered(
		common.HexToAddress(os.Getenv("TRANSFER_ID")),
		types.Token{
			Type:    types.TokenTypeERC20,
			ChainId: types.ChainIdBase,
			Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
		},
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clRecovered)
	url, err := clRecovered.GenerateClaimUrl(
		getRandomBytes,
		signTypedData,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(url)
}
