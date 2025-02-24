package main

import (
	"crypto/rand"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"log"
	"math/big"
	"os"
)

func signTypedData(typedData apitypes.TypedData) ([]byte, error) {
	pkHex := os.Getenv("PRIVATE_KEY") // Replace with your method to get the private key
	privKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Hash the typed data
	typedDataHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, fmt.Errorf("failed to hash typed data: %w", err)
	}

	// Sign the hash with the private key
	signature, err := crypto.Sign(typedDataHash[:], privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the typed data: %w", err)
	}

	// Adjust the signature's "v" value to align with Ethereum's standards
	signature[64] += 27 // Ethereum appends 27 to v (recovery id) to construct a full signature

	return signature, nil
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

	claimLink, err := sdk.CreateClaimLink(
		types.Token{
			Type:    types.TokenTypeNative,
			ChainId: types.ChainIdBase,
		},
		big.NewInt(1000000000),
		common.HexToAddress("0x3A205ECf286bBe11460638aCe47D501A53fB91C0"),
		big.NewInt(1000000000),
	)
	if err != nil {
		log.Fatalln(err)
	}

	link, transferId, err := claimLink.GenerateClaimUrl(signTypedData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(link, transferId)
}
