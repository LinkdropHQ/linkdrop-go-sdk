package main

import (
	"crypto/rand"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/helpers"
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
		os.Getenv("LINKDROP_API_KEY"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Native
	clNative, err := sdk.ClaimLink(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeNative,
				ChainId: types.ChainIdBase,
			},
			Sender:     common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
			Amount:     big.NewInt(250000000000000000),
			Expiration: 1775195026,
		},
		getRandomBytes,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clNative)

	// ERC20 claim link with pre-generated linkKey provided
	// Any method can be used to generate ecdsa.PrivateKey
	linkKey, err := helpers.PrivateKey(getRandomBytes)
	if err != nil {
		log.Fatalln(err)
	}
	clERC20WithLinkKey, err := sdk.ClaimLinkWithLinkKey(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdBase,
				Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
			},
			Sender:     common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
			Amount:     big.NewInt(1000000000),
			Expiration: 1775195026,
		},
		*linkKey,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clERC20WithLinkKey)

	// ERC20 claim link with transferId
	transferId := common.HexToAddress("0xcc06431Bcb7E5BDf5632705db6Eb4e98123e3e78")
	clERC20WithTransferId, err := sdk.ClaimLinkWithTransferId(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdBase,
				Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
			},
			Sender:     common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
			Amount:     big.NewInt(1000000000),
			Expiration: 1775195026,
		},
		transferId,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clERC20WithTransferId)
}
