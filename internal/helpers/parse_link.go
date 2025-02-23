package helpers

import (
	"crypto/ecdsa"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
)

func ParseLink(
	link string,
	decodedLink *types.Link,
) (senderSig string, linkKey *ecdsa.PrivateKey, err error) {
	if decodedLink == nil {
		decodedLink, err = DecodeLink(link)
		if err != nil {
			return "", nil, err
		}
	}
	return decodedLink.SenderSig, decodedLink.LinkKey, nil
}
