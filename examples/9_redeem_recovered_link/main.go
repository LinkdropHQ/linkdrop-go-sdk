package main

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"os"
)

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
		linkdrop.WithCoinbaseWalletProductionDefaults(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	link, err := sdk.GetClaimLink(os.Getenv("LINKDROP_RECOVERED_LINK"))
	if err != nil {
		log.Fatalln(err)
	}

	txHash, err := link.Redeem(common.HexToAddress(os.Getenv("RECEIVER_ADDRESS")))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Redeem transaction hash:", txHash)
}
