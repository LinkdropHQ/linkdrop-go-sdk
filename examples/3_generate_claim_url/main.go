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
		types.DeploymentCBW,
		getRandomBytes,
		linkdrop.WithApiKey(os.Getenv("LINKDROP_API_KEY")),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// ERC20
	clERC20, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeERC20,
			ChainId: types.ChainIdBase,
			Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
		},
		big.NewInt(100000),
		common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
		big.NewInt(1695985897077),
	)
	if err != nil {
		log.Fatalln(err)
	}
	url, _, err := clERC20.GenerateClaimUrl(signTypedData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(url)
	log.Println("This link is valid but can't be claimed yet. See 7_deposit for the example of depositing assets.")
}
