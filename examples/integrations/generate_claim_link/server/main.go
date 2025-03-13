package main

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"os"
	"strings"
)

func main() {
	linkKey, err := crypto.HexToECDSA(strings.TrimPrefix(os.Args[1], "0x"))

	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// ERC20
	clERC20, err := sdk.ClaimLinkWithLinkKey(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdBase,
				Address: common.HexToAddress("0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"),
			},
			Sender:     common.HexToAddress("0x5659A8557FdBA11AA04cfCfcc59EeF9FA412A7dD"),
			Amount:     big.NewInt(100000),
			Expiration: 1773159165,
		},
		*linkKey,
	)
	if err != nil {
		log.Fatalln(err)
	}
	url, err := clERC20.GenerateClaimUrl()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print("Claim Link: ", url)
}
