package main

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/LinkdropHQ/linkdrop-go-sdk/utils"
	"log"
	"os"
)

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		types.DeploymentCBW,
		utils.GetRandomBytes,
		linkdrop.WithApiKey(os.Getenv("LINKDROP_API_KEY")),
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
