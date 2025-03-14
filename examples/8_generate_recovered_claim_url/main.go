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
	)
	if err != nil {
		log.Fatalln(err)
	}

	// ERC20
	// https://p2p.linkdrop.io/#/code?k=H4RpFzgphRTGzdiVTy84EYzDTXWeWSY8VQwukQD94NU2&c=8453&v=3&src=p2p
	clERC20, err := sdk.ClaimLinkWithTransferId(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdBase,
				Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
			},
		},
		common.HexToAddress("0x2021f4B02E2B59b89E5D2c8028Da8DC4d1Ef1Fb8"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(clERC20)
	url, err := clERC20.GenerateRecoveredClaimUrl(
		getRandomBytes,
		signTypedData,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(url)
}
