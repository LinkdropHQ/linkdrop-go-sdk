package main

import (
	"encoding/json"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"os"
)

type Payload struct {
	TransferId common.Address `json:"transferId"`
	LinkKeyId  common.Address `json:"linkKeyId"`
	ClaimLink  struct {
		Token struct {
			Type    types.TokenType `json:"type"`
			ChainId types.ChainId   `json:"chainId"`
			Address common.Address  `json:"address"`
		} `json:"token"`
	} `json:"claimLink"`
}

func main() {
	command := os.Args[1]
	var payload Payload
	err := json.Unmarshal([]byte(os.Args[2]), &payload)
	if err != nil {
		log.Fatalln(err)
	}

	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	claimLink, err := sdk.ClaimLinkRecovered(
		payload.TransferId,
		types.Token{
			Type:    payload.ClaimLink.Token.Type,
			ChainId: payload.ClaimLink.Token.ChainId,
			Address: payload.ClaimLink.Token.Address,
		},
		nil, nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	switch command {
	case "getRecoveredLinkTypedData":
		params, err := claimLink.GetTypedData(payload.LinkKeyId)
		if err != nil {
			log.Fatalln(err)
		}
		resp, _ := json.Marshal(params)
		fmt.Println(string(resp))
	}
}
