package main

import (
	"crypto/rand"
	"github.com/LinkdropHQ/linkdrop-go-sdk"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
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
	client, err := linkdrop.Init(
		"https://p2p.linkdrop.io",
		types.DeploymentCBW,
		getRandomBytes,
		linkdrop.WithApiKey(os.Getenv("LINKDROP_API_KEY")),
	)
	if err != nil {
		panic(err)
	}
	log.Println(client)
}
