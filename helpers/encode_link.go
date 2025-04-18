package helpers

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/mr-tron/base58"
	"strconv"
)

func EncodeLink(claimHost string, link types.Link) string {
	// Encode LinkKey and TransferId to Base58
	linkKey := base58.Encode(link.LinkKey.D.Bytes())
	transferId := base58.Encode(link.TransferId.Bytes())

	// Handle optional encryption key
	var encryptionKey string
	if link.Message != nil && link.Message.LinkKey != "" {
		encryptionKey = "&m=" + string(link.Message.LinkKey)
	}
	chainId := strconv.Itoa(int(link.ChainId))

	// Handle optional SenderSig
	if link.SenderSignature != nil {
		sig := base58.Encode(link.SenderSignature)
		return fmt.Sprintf("%s/#/code?k=%s&sg=%s&i=%s&c=%s&v=3&sgl=%d&src=p2p%s",
			claimHost, linkKey, sig, transferId, chainId, len(link.SenderSignature), encryptionKey)
	}

	// If SenderSig is not provided
	return fmt.Sprintf("%s/#/code?k=%s&c=%s&v=3&src=p2p%s", claimHost, linkKey, chainId, encryptionKey)
}
