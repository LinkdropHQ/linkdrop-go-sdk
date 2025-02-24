package main

import (
	"crypto/rand"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"log"
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
	link, err := sdk.GetClaimLink(
		os.Getenv("TEST_CLAIM_LINK"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	txHash, err := link.Redeem(common.HexToAddress(os.Getenv("RECEIVER_ADDRESS")))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Redeem transaction hash:", txHash)
}
