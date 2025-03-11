package types

import (
	"crypto/sha256"
	"github.com/mr-tron/base58"
)

type EncryptedMessage struct {
	Data    []byte
	LinkKey MessageLinkKey
}

type MessageInitialKey [32]byte

func (mki *MessageInitialKey) LinkKey(
	length uint16,
) MessageLinkKey {
	initialKeyEncoded := base58.Encode(mki[:])
	if len(initialKeyEncoded) > int(length) {
		initialKeyEncoded = initialKeyEncoded[:length]
	}
	return MessageLinkKey(initialKeyEncoded)
}

type MessageLinkKey string

func (mlk *MessageLinkKey) MessageEncryptionKey() (encryptionKey MessageEncryptionKey, err error) {
	initialKeyTrimmed, err := base58.Decode(string(*mlk))
	if err != nil {
		return
	}
	encryptionKey = sha256.Sum256(initialKeyTrimmed)
	return
}

type MessageEncryptionKey [32]byte
