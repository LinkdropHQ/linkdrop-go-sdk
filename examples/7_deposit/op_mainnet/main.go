package main

import (
	"crypto/rand"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/LinkdropHQ/linkdrop-go-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func sendTransaction(chainId *big.Int, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}
	return utils.SendTransaction(chainId, to, value, data, client, privateKey)
}

func main() {
	sdk, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		os.Getenv("LINKDROP_API_KEY"),
		linkdrop.WithProductionDefaults(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// An example with a custom Escrow Address
	// (can be skipped to use the default one for the Token chain set in ClaimLink)
	escrowAddress, _, err := utils.EscrowAddressByChain(types.ChainIdOptimism)
	if err != nil {
		log.Fatalln(err)
	}

	// ERC20
	clERC20, err := sdk.ClaimLink(
		linkdrop.ClaimLinkCreationParams{
			Token: types.Token{
				Type:    types.TokenTypeERC20,
				ChainId: types.ChainIdOptimism,
				Address: common.HexToAddress("0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85"),
			},
			Sender:        common.HexToAddress(os.Getenv("SENDER_ADDRESS")),
			Amount:        big.NewInt(100000),
			Expiration:    1773234550,
			EscrowAddress: &escrowAddress, // Custom Escrow Address
		},
		getRandomBytes,
	)
	if err != nil {
		log.Fatalln(err)
	}

	url, err := clERC20.ClaimUrl()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(url) // The link is valid, but can't be claimed since assets are were deposited

	txHash, err := clERC20.Deposit(sendTransaction)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("TX Hash: ", txHash)
}
