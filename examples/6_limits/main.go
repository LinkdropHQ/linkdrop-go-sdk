package main

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"log"
	"os"
)

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	limits, err := sdk.GetLimits(
		types.Token{
			Type:    types.TokenTypeNative,
			ChainId: types.ChainIdBase,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(limits)
}
