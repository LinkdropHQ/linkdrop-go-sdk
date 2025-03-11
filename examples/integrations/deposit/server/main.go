package main

import (
	"encoding/json"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"os"
)

type Payload struct {
	Command    string         `json:"command"`
	TransferId common.Address `json:"transferId"`
	ClaimLink  struct {
		Token struct {
			Type    types.TokenType `json:"type"`
			ChainId types.ChainId   `json:"chainId"`
			Address common.Address  `json:"address"`
		} `json:"token"`
		Sender     common.Address `json:"sender"`
		Amount     string         `json:"amount"`
		Expiration int64          `json:"expiration"`
	} `json:"claimLink"`
}

func main() {
	var payload Payload
	err := json.Unmarshal([]byte(os.Args[1]), &payload)
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

	//// ERC20
	amount, _ := new(big.Int).SetString(payload.ClaimLink.Amount, 10)
	claimLink, err := sdk.ClaimLinkWithTransferId(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    payload.ClaimLink.Token.Type,
				ChainId: payload.ClaimLink.Token.ChainId,
				Address: payload.ClaimLink.Token.Address,
			},
			Sender:     payload.ClaimLink.Sender,
			Amount:     amount,
			Expiration: payload.ClaimLink.Expiration,
		},
		payload.TransferId,
	)
	if err != nil {
		log.Fatalln(err)
	}

	switch payload.Command {
	case "getDepositParams":
		params, err := claimLink.GetDepositParams()
		if err != nil {
			log.Fatalln(err)
		}
		resp, _ := json.Marshal(params)
		fmt.Println(string(resp))
	}
}
