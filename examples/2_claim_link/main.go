package main

import (
	"crypto/rand"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		types.DeploymentCBW,
		getRandomBytes,
		linkdrop.WithApiKey(os.Getenv("LINKDROP_API_KEY")),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Native
	clNative, err := sdk.CreateClaimLink(
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
	log.Println(sdk, clNative)

	// ERC20
	clERC20, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeERC20,
			ChainId: types.ChainIdBase,
			Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
		},
		big.NewInt(1000000000),
		common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
		big.NewInt(1000000000),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sdk, clERC20)

	// ERC721
	clERC721, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeERC721,
			ChainId: types.ChainIdBase,
			Address: common.HexToAddress("0x3319197b0d0f8ccd1087f2d2e47a8fb7c0710171"),
			Id:      big.NewInt(5225),
		},
		big.NewInt(1),
		common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
		big.NewInt(1000000000),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sdk, clERC721)
}
