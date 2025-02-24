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
	claimLink, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeNative,
			ChainId: types.ChainIdBase,
		},
		big.NewInt(1000000000),
		common.HexToAddress("0x3A205ECf286bBe11460638aCe47D501A53fB91C0"),
		big.NewInt(1000000000),
	)
	log.Println(sdk, claimLink)
}
