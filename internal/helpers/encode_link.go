package helpers

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/mr-tron/base58"
)

func EncodeLink(claimHost string, link types.Link) string {
	// Encode LinkKey and TransferId to Base58
	linkKey := base58.Encode(link.LinkKey.D.Bytes())
	transferId := base58.Encode(link.TransferId.Bytes())

	// Handle optional encryption key
	var encryptionKey string
	if link.EncryptionKey != nil {
		// TODO 0x prefix?
		encryptionKey = fmt.Sprintf("&m=0x%x", link.EncryptionKey)
	}

	// Handle optional SenderSig
	if link.SenderSig != "" {
		sigLength := (len(link.SenderSig) - 2) / 2
		sig := base58.Encode([]byte(link.SenderSig))
		return fmt.Sprintf("%s/#/code?k=%s&sg=%s&i=%s&c=%s&v=3&sgl=%d&src=p2p%s",
			claimHost, linkKey, sig, transferId, link.ChainId, sigLength, encryptionKey)
	}

	// If SenderSig is not provided
	return fmt.Sprintf("%s/#/code?k=%s&c=%s&v=3&src=p2p%s", claimHost, linkKey, link.ChainId, encryptionKey)
}
